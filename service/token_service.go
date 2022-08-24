package service

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const secretKey = "eyJhbGciOiJSUzI1NiIsInR"

type Claims struct {
	Username string
	jwt.StandardClaims
}

// GenerateToken 生成Token值
func GenerateToken(userName string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(86400 * 14 * time.Second)
	claims := Claims{
		Username: userName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))
	return token, err
}

// token: "eyJhbGciO...解析token"
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	fmt.Println("-=====")

	return nil, err
}
