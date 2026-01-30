package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/Saikatdeb12/TodoApp2/internal/database"
	"github.com/google/uuid"
)

const UserIDkey string = "user_id"

func Auth(next http.Handler) http.Handler{
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
		token := r.Header.Get("Authorization")
		if token == ""{
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		sessionID, err := uuid.Parse(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		var userID uuid.UUID
		var expires time.Time

		query := `
			SELECT user_id, expires_at
			FROM sessions
			WHERE id=$1
		`

		err = database.DB.QueryRow(query, sessionID).Scan(&userID, &expires)
		if err != nil || expires.Before(time.Now()){
			http.Error(w, "Session expired", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDkey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
