package server

import (
	"errors"
	"time"

	"referal/internal/config"

	"github.com/dgrijalva/jwt-go"
)

type Token struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func MakeToken(email string) (string, error) {
	expires := time.Now().Add(24 * time.Hour)
	claims := &Token{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expires.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenS, err := token.SignedString(config.SecretKey)

	if err != nil {
		return "", errors.New("ошибка создания токена")
	}

	return tokenS, nil
}

func DecodeJWT(tokenStr string) (string, error) {
	claims := &Token{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return config.SecretKey, nil
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
