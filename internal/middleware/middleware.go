package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/Saikatdeb12/TodoApp2/internal/database"
	"github.com/Saikatdeb12/TodoApp2/internal/models"
	"github.com/google/uuid"
)

const RequestContextKey string = "request_context"

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		var sessionExpiresAt time.Time

		// err = database.DB.QueryRow(query, sessionID).Scan(&userID, &sessionExpiresAt)
		userID, err := database.ValidateUserSession(sessionID.String())

		if err != nil || sessionExpiresAt.Before(time.Now()) {
			http.Error(w, "Session expired", http.StatusUnauthorized)
			return
		}

		requestContext := models.RequestContext{
			UserID:    userID,
			SessionID: sessionID,
		}

		ctx := context.WithValue(r.Context(), RequestContextKey, requestContext)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
