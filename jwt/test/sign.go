package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func Sign(secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  "1234567890",
		"name": "zomasec",
		"iat":  time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	return token.SignedString([]byte(secret))
}


func main() {
	secret := "secret123"
	token, err := Sign(secret)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(token)
}