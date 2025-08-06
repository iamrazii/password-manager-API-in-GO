package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(ID int, username string) (string, error) {

	// creating payload
	claims := jwt.MapClaims{
		"id":       ID,
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(), // expires in 1 day
	}

	// assigning claims and defining signing algorithm
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// getting secret key for jwt token from .env file
	JWT_Secret := []byte(os.Getenv("JWT_SECRET"))

	fmt.Printf("utils has %s\n", JWT_Secret)

	// returning our JWT token after signing it with our JWT_Secret key
	return token.SignedString(JWT_Secret)

}
