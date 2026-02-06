package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/Saikatdeb12/TodoApp2/internal/models"
	"github.com/Saikatdeb12/TodoApp2/internal/utils"
	"github.com/golang-jwt/jwt"
)

type ContextKeys struct{}

var RequestContextKey = ContextKeys{}

func UserContext(r *http.Request) *models.RequestContext {
	user, _ := r.Context().Value(RequestContextKey).(*models.RequestContext)
	return user
}

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { // HOF
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.RespondError(w, http.StatusUnauthorized, nil, "missing authorization header")
			return
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			utils.RespondError(w, http.StatusUnauthorized, nil, "invalid authorization format")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, bearerPrefix)

		token, parseErr := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, err := token.Method.(*jwt.SigningMethodHMAC); !err {
				return nil, errors.New("invalid signing method")
			}
			return []byte(utils.SecretKey), nil
		})

		if parseErr != nil || !token.Valid {
			utils.RespondError(w, http.StatusUnauthorized, parseErr, "invalid token")
			return
		}

		claimValues, err := token.Claims.(jwt.MapClaims)

		if !err {
			utils.RespondError(w, http.StatusUnauthorized, nil, "invalid token claims")
			return
		}

		userID, err := claimValues["user_id"].(string)
		if !err {
			utils.RespondError(w, http.StatusUnauthorized, nil, "invalid user id")
		}

		sessionID, err := claimValues["session_id"].(string)
		if !err {
			utils.RespondError(w, http.StatusUnauthorized, nil, "invalid session id")
		}

		requestContext := models.RequestContext{
			UserID:    userID,
			SessionID: sessionID,
		}

		ctx := context.WithValue(r.Context(), RequestContextKey, requestContext)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
