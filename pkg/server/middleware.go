package server

import (
	"fmt"
	"net/http"
	"referal/pkg/log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type writer struct {
	http.ResponseWriter
	statusCode int
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		log.Logger.Info(fmt.Sprintf("Request %s %s", r.Method, r.URL.Path))
		
		wrappedWriter := &writer{w, http.StatusOK}
		next.ServeHTTP(wrappedWriter, r)

		log.Logger.Info(fmt.Sprintf("Completed %s %s with %d in %v", r.Method, r.URL.Path, wrappedWriter.statusCode, time.Since(start)))
	})
}

func CheckJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("JWT")

		if err != nil {
			if err == http.ErrNoCookie {
				http.Redirect(w, r, "/signin", http.StatusSeeOther)
				return
			}

			http.Error(w, "Ошибка при получении куки", http.StatusBadRequest)
			return
		}

		tokenS := cookie.Value
		claims := &Token{}

		token, err := jwt.ParseWithClaims(tokenS, claims, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Redirect(w, r, "/signin", http.StatusSeeOther)
				return
			}
			
			http.Error(w, "Ошибка при создании токена", http.StatusBadRequest)
			return
		}

		if !token.Valid {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}