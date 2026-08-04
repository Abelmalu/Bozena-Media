package pkg

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	accessSecret  = []byte("ACCESS_SECRET_CHANGE_ME")
	refreshSecret = []byte("REFRESH_SECRET_CHANGE_ME")
)

const RefreshTokenDuration = 24 * 30 * time.Hour

func GenerateAcessToken(userID int,role string) (string, error) {

	claims := jwt.MapClaims{
		"user_id": userID,
		"type":    "access",
		"userRole":role,
		"exp":     time.Now().Add(time.Minute * 15).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(accessSecret)
}

// generates refresh tokens with refresh token secret keys
func GenerateRefreshToken(userID int) (string, error, time.Time) {
	expiresAt := time.Now().Add(RefreshTokenDuration)

	claims := jwt.MapClaims{
		"user_id": userID,
		"type":    "refresh",
		"exp":     expiresAt,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(refreshSecret)

	return signedToken, err, expiresAt

}

func ValidateAccessToken(tokenStr string) (jwt.MapClaims, error) {

	return ValidateToken(tokenStr, accessSecret, "access")
}

func ValidateRefreshToken(tokenStr string) (jwt.MapClaims, error) {
	
	return ValidateToken(tokenStr, refreshSecret, "refresh")
}

// validate both access and refresh tokens
func ValidateToken(tokenStr string, secret []byte, expectedType string) (jwt.MapClaims, error) {
	// 1. Parse the token
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {

		// Verify the signing method is what you expect (HMAC)
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		// This returns the specific error (e.g., "token is expired")
		return nil, fmt.Errorf("parsing token failed: %w", err)
	}

	// Extract Claims
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, errors.New("invalid token claims structure")
	}

	// 5. Validate the "type" claim (Access vs Refresh)
	// We safely cast to string to avoid potential interface{} comparison issues
	typeVal, ok := claims["type"].(string)
	if !ok || typeVal != expectedType {
		return nil, fmt.Errorf("invalid token type: expected %s, got %v", expectedType, claims["type"])
	}

	return claims, nil
}

// HashToken hashes refresh tokens using SHA256
func HashToken(token string) string {

	hasher := sha256.New()
	hasher.Write([]byte(token))
	return hex.EncodeToString(hasher.Sum(nil))

}
