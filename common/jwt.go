package common

import (
	"github.com/dgrijalva/jwt-go"
	"magicMakeup/entities"
	"time"
)

var jwtKey = []byte("a_secret_key")

type Claims struct {
	UserID int
	jwt.StandardClaims
}

func ReleaseToken(user *entities.User) (string, error) {
	expireTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserID:         user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			IssuedAt: time.Now().Unix(),
			Issuer: "Magic Make Up",
			Subject: "User Token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, nil, err
	}
	return token, claims, nil
}
