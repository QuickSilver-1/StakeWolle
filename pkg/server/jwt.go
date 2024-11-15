package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var secretKey = []byte("M8axEo25vLElQ8n85CvmFRmNrFWt0YQq")

type Token struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func MakeToken(email string, w http.ResponseWriter) error {
	token, err := newToken(email)

	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name: "JWT",
		Value: token,
		Expires: time.Now().Add(100 * time.Hour),
	})

	return nil
}

func newToken(email string) (string, error) {
	expires := time.Now().Add(100 * time.Hour)
	claims := &Token{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expires.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenS, err := token.SignedString(secretKey)

	if err != nil {
		return "", errors.New("ошибка создания токена")
	}

	return tokenS, nil
}

func DecodeJWT(tokenStr string) (string, error) {
	claims := &Token{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", errors.New("неверная форма токена")
		}
		
		return "", err
	}
	
	if !token.Valid {
		return "", errors.New("сломанный токен")
	}

	return claims.Email, nil 
}
