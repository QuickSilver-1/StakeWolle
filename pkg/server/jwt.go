package server

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var secretKey = []byte("M8axEo25vLElQ8n85CvmFRmNrFWt0YQq")

type Token struct {
	email string `json:"email"`
	jwt.StandardClaims
}

func NewToken(email string) (string, error) {
	expires := time.Now().Add(100 * time.Hour)
	claims := &Token{
		email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expires.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenS, err := token.SignedString(secretKey)

	if err != nil {
		return "", err
	}

	return tokenS, err
}
