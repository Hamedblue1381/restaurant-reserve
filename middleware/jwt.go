package middleware

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	UserId uint `json:"id"`
	Email string `json:"email"`
	Role string `json:"role"`
	jwt.StandardClaims
}

// Generate Token for a given email ✨
func GenerateToken(email string, userId uint, role string) (string, error) {

	exprTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Email: email,
		UserId: userId,
		Role: role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exprTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	return tokenString, err
}

func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.NewValidationError("Invalid Token", jwt.ValidationErrorMalformed)
	}

	return claims, nil
}