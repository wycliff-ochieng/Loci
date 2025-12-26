package middleware

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/wycliff-ochieng/pkg/auth"
	//"golang.org/x/net/context"
)

type ContextKey string

var (
	UserIDKey ContextKey = "userid"
)

func AuthenticationMiddleware(jwtSecret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "authHeader cant be empty,Authorization header is required", http.StatusFailedDependency)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader || tokenString == "" {
				http.Error(w, "invalid token", http.StatusFailedDependency)
				return
			}

			//parse and validate token using shared Claims
			token, err := jwt.ParseWithClaims(tokenString, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method, %s", token.Header["alg"])
				}
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				log.Printf("Token Issue: %s", err)
				http.Error(w, "issue with token", http.StatusExpectationFailed)
				return
			}
			//type assertion
			if claims, ok := token.Claims.(*auth.Claims); ok && claims != nil {
				//populate context
				ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				http.Error(w, "could not parse the claims into context", http.StatusFailedDependency)
			}
		})
	}
}

//get USerID from context
/*
func GetUserUUIDFromContext(ctx context.Context) (uuid.UUID, error) {
	userUUIDStr, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return uuid.Nil, errors.New("user UUID not found in context")
	}
	if userUUIDStr == "" {
		return uuid.Nil, errors.New("user UUID in context is empty")
	}

	parsedUUID, err := uuid.Parse(userUUIDStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid UUID format in context: %w", err)
	}

	return parsedUUID, nil
}
*/
func GetUserUUIDFromContext(ctx context.Context) (uuid.UUID, error) {
	userUUID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("userID not found in context")
	}

	return userUUID, nil
}

//RBAC if required
