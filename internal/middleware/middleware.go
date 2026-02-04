package middlewares

import (
	"context"
	"net/http"

	"github.com/Saikatdeb12/TodoApp2/internal/database"
	"github.com/Saikatdeb12/TodoApp2/internal/models"
	"github.com/google/uuid"
)

type ContextKeys struct{}

var RequestContextKey = ContextKeys{}

func UserContext(r *http.Request) *models.RequestContext {
	user, _ := r.Context().Value(RequestContextKey).(*models.RequestContext)
	return user
}

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { // HOC
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		sessionID, err := uuid.Parse(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		userID, err := database.ValidateUserSession(sessionID)

		// if err != nil || sessionExpiresAt.Before(time.Now()) {
		// 	http.Error(w, "Session expired", http.StatusUnauthorized)
		// 	return
		// }

		requestContext := models.RequestContext{
			UserID:    userID,
			SessionID: sessionID,
		}

		ctx := context.WithValue(r.Context(), RequestContextKey, requestContext)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
