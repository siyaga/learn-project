package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))
var jwtRefreshKey = []byte(os.Getenv("JWT_SECRET_REFRESH"))

// Claims struct untuk token JWT
type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// Generate JWT Access Token (Berlaku 1 Jam)
func GenerateToken(email string) (string, string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Berlaku 1 hari

	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	return signedToken, "1 day", nil
}
// Generate Refresh Token (Berlaku 7 Hari)
func GenerateRefreshToken(email string) (string, error) {
	expirationTime := time.Now().Add(time.Hour * 24 * 7) // Berlaku 7 hari
	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtRefreshKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// Validasi Token
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}
