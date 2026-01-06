package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(userID int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = userID
	claims["exp"] = time.Now().Add(15 * time.Minute).Unix()

	secret := os.Getenv("ACCESS_TOKEN_SECRET")
	tokenString, err := token.SignedString([]byte(secret)) // Assina o token
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
