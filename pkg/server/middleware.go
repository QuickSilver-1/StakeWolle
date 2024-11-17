package server

import (
	"fmt"
	"net/http"
	"referal/internal/config"
	"referal/pkg/log"
	"strings"
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
		tokenS := r.Header.Get("Authorization")

		if tokenS == "" {
			AnswerHandler(w, 401, "Требуется авторизация")
			return
		}

		tokenS = strings.Split(tokenS, " ")[1]
		claims := &Token{}

		token, err := jwt.ParseWithClaims(tokenS, claims, func(token *jwt.Token) (interface{}, error) {
			return config.SecretKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				AnswerHandler(w, 400, "Ошибка сигнатуры токена")
				return
			}
			
			AnswerHandler(w, 500, "Неверный токен")
			return
		}

		if !token.Valid {
			AnswerHandler(w, 401, "Необходима авторизация")
			return
		}

		next.ServeHTTP(w, r)
	})
}