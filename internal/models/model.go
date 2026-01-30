package models
import (
	"time"
	"github.com/google/uuid"
)

type User struct {
	UserID uuid.UUID `json:"id" db:"user_id"`
	Name string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
	Password string `json:"-" db:"password"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	ArchivedAt *time.Time `json:"archivedAt" db:"archived_at"`
}

type Todo struct {
	TodoID uuid.UUID `json:"id" db:"id"`
	UserID uuid.UUID `json:"userID" db:"user_id"`
	Title string `json:"title" db:"title"`
	Body string `json:"body" db:"body"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	ValidTill time.Time `json:"validTill" db:"valid_till"`
	Complete  bool  `json:"complete" db:"complete"`
}

type Session struct{
	SessionID uuid.UUID `json:"-" db:"session_id"`
	UserID uuid.UUID `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
}

type CreateTodoRequest struct {
	Title string `json:"title"`
	Body string `json:"body"`
	ValidTill time.Time `json:"validTill"`
}

type UpdateTodoRequest struct{
	Title *string `json:"title"`
	Body *string `json:"body"`
	Complete *bool `json:"complete"`
	ValidTill *time.Time `json:"valid_till"`
}

type UpcomingTodosRequest struct{
	Title *string `json:"title"`
	Body *string `json:"body"`
	Complete *bool `json:"complete"`
	ValidTill *time.Time `json:"valid_till"`
}

type RegisterRequest struct{
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct{
	Email string `json:"email"`
	Password string `json:"password"`
}

