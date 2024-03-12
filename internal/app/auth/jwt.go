package auth

import (
	"fmt"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/logger"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"time"
)

const (
	secretKey = "GopherSecretKey"
	tokenExp  = time.Hour * 24
)

var authLog = logger.Log.With(zap.String("cat", "AUTH"))

func GenerateToken(userID string) (string, error) {
	authLog.Debug(fmt.Sprintf("creating new token for sub = %s", userID))
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("signedString. %w", err)
	}
	return tokenString, nil
}

func ValidateToken(tokenString string) (*jwt.RegisteredClaims, bool) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		authLog.Debug("parsing jwt with claims", zap.String("token", tokenString),
			zap.Error(err))
		return nil, false
	}
	return claims, token.Valid
}
