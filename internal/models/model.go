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

type RegisterRequest struct{
	Name string `json:"name" validate:"required,min=2,max=50"`
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=8,lte=30"`
}

type LoginRequest struct{
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CreateTodoRequest struct {
	Title      string    `json:"title" validate:"required,min=3,max=100"`
	Body string `json:"body" validate:"required"`
	ValidTill time.Time `json:"validTill" validate:"required"`
}

type UpdateTodoRequest struct{
	Title *string `json:"title" validate:"omitempty,min=3,max=100"`
	Body *string `json:"body" validate:"omitempty"`
	Complete *bool `json:"complete"`
	ValidTill *time.Time `json:"valid_till"`
}

type UpcomingTodosRequest struct{
	Title *string `json:"title"`
	Body *string `json:"body"`
	Complete *bool `json:"complete"`
	ValidTill *time.Time `json:"valid_till"`
}

