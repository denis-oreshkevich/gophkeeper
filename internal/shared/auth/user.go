package auth

import (
	"context"
	"errors"
)

type UserIDKey struct{}

func GetUserID(ctx context.Context) (string, error) {
	value := ctx.Value(UserIDKey{})
	if value == nil {
		return "", errors.New("userID is not present in context")
	}
	userID, ok := value.(string)
	if !ok {
		return "", errors.New("userID is not string")
	}
	return userID, nil
}
