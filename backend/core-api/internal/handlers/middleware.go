package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type contextKey string

const userContextKey = contextKey("user")

var jwtSecret = []byte("bardzo-sekretny-klucz-do-zmiany-pozniej")

func (h *Handlers) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized: missing authorization header", http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			http.Error(w, "Unauthorized: malformed authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := headerParts[1]

		claims := &jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized: invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, claims.Subject)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handlers) getUserIDFromContext(r *http.Request) (uuid.UUID, error) {
	userIDStr, ok := r.Context().Value(userContextKey).(string)
	if !ok {
		return uuid.Nil, errors.New("unable to retrieve user ID from context")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}
