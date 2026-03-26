// Crack This :
package main

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)


func main() {

	tokenStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3NzM3Mjk1MjIsIm5hbWUiOiJ6b21hc2VjIiwic3ViIjoiMTIzNDU2Nzg5MCJ9.LKylehXrNzgjwdsIj5IsEKjlYax4EWYIXPc7TFOwsGQ"
	parser := jwt.NewParser()

	token, _, err := parser.ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		fmt.Println("Error parsing token:", err)
		return
	}

	Algo := token.Header["alg"] 
	fmt.Println("Algorithm:", Algo)

	method, ok := token.Method.(*jwt.SigningMethodHMAC) 
	if !ok {
		fmt.Println("Method is not brute-forceable")
		return
	}

	secretsList, err := loadWordlist("/home/zomasec/automation-course/vulns/jwt/secrets.list")
	if err != nil {
		fmt.Println("Error loading wordlist:", err)
		return
	}

	signingString, _ := token.SigningString()
	for _, secret := range secretsList {
		if err := method.Verify(signingString, token.Signature, []byte(secret)); err == nil {
			fmt.Println("Secret found:", secret)
			break
		}
	}

	
}
