package middleware

import (
	"context"
	"errors"
	"net/http"

	dbhelper "github.com/Saikatdeb12/TodoApp2/internal/database/dbHelper"
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
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		userID, err := dbhelper.ValidateUserSession(tokenString)
		if err != nil {
			utils.RespondError(w, http.StatusUnauthorized, err, "session expired")
			return
		}

		token, parseErr := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return []byte(utils.SecretKey), nil
		})

		if parseErr != nil || !token.Valid {
			utils.RespondError(w, http.StatusUnauthorized, parseErr, "invalid token")
			return
		}

		claimValues, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			utils.RespondError(w, http.StatusUnauthorized, nil, "invalid token claims")
			return
		}

		requestContext := models.RequestContext{
			UserID:    userID,
			SessionID: claimValues["sessionId"].(string),
		}

		ctx := context.WithValue(r.Context(), RequestContextKey, requestContext)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
