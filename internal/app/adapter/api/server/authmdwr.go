package server

import (
	"context"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/auth"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/logger"
	"go.uber.org/zap"
	"net/http"
)

const AuthorizationHeaderName = "Authorization"

var whiteList = map[string]struct{}{
	"/api/user/register": {},
	"/api/user/login":    {},
}

var log = logger.Log.With(zap.String("cat", "AUTH"))

func Auth(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		uri := r.RequestURI
		_, ok := whiteList[uri]
		if ok {
			next.ServeHTTP(w, r)
		}

		tokenString := r.Header.Get(AuthorizationHeaderName)
		if tokenString == "" {
			log.Debug(AuthorizationHeaderName + " header not found")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		claims, isValid := auth.ValidateToken(tokenString)
		if !isValid {
			log.Debug("token is not valid")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, auth.UserIDKey{}, claims.Subject)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(f)
}
