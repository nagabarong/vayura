package pkg

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret string
var jwtExpiresIn time.Duration

// SetJWTSecret sets the JWT secret key
func SetJWTSecret(secret string) {
	jwtSecret = secret
}

// SetJWTExpiration sets the JWT expiration duration
func SetJWTExpiration(d time.Duration) {
	jwtExpiresIn = d
}

// GenerateJWT generates a JWT token for the user
func GenerateJWT(userID uint, email string) (string, error) {
	if jwtSecret == "" {
		return "", errors.New("JWT secret not configured")
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(getExpiration()).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func getExpiration() time.Duration {
	if jwtExpiresIn <= 0 {
		return 72 * time.Hour
	}
	return jwtExpiresIn
}

// VerifyJWT verifies and parses a JWT token
func VerifyJWT(tokenString string) (*jwt.Token, error) {
	if jwtSecret == "" {
		return nil, errors.New("JWT secret not configured")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})

	return token, err
}

// ExtractClaims extracts claims from a JWT token
func ExtractClaims(token *jwt.Token) (jwt.MapClaims, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}
