package server

import (
	"errors"
	"fmt"
	"time"

	"referal/internal/config"
	"referal/pkg/log"

	"github.com/dgrijalva/jwt-go"
)

type Token struct {
    Email string `json:"email"`
    jwt.StandardClaims
}

// MakeToken создает JWT токен для данного email
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
        log.Logger.Error(fmt.Sprintf("Ошибка создания токена для email %s: %v", email, err))
        return "", errors.New("ошибка создания токена")
    }

    log.Logger.Info(fmt.Sprintf("Токен успешно создан для email %s", email))
    return tokenS, nil
}

// DecodeJWT декодирует JWT токен и возвращает email
func DecodeJWT(tokenStr string) (string, error) {
    claims := &Token{}
    token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
        return config.SecretKey, nil
    })

    if err != nil {
        if err == jwt.ErrSignatureInvalid {
            log.Logger.Error(fmt.Sprintf("Неверная форма токена: %v", err))
            return "", errors.New("неверная форма токена")
        }
        
        log.Logger.Error(fmt.Sprintf("Ошибка при декодировании JWT: %v", err))
        return "", err
    }
    
    if !token.Valid {
        log.Logger.Warn("Сломанный токен")
        return "", errors.New("сломанный токен")
    }

    log.Logger.Info(fmt.Sprintf("Токен успешно декодирован для email %s", claims.Email))
    return claims.Email, nil 
}
