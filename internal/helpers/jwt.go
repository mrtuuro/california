package helpers

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Email": email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	return tokenStr, err
}
