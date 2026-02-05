package models

import (
	"time"
)

type User struct {
	UserID     string     `json:"id" db:"user_id"`
	Name       string     `json:"name" db:"name"`
	Email      string     `json:"email" db:"email"`
	Password   string     `json:"-" db:"password"`
	CreatedAt  time.Time  `json:"createdAt" db:"created_at"`
	ArchivedAt *time.Time `json:"archivedAt" db:"archived_at"`
}

type Todo struct {
	TodoID    string    `json:"id" db:"id"`
	UserID    string    `json:"userID" db:"user_id"`
	Title     string    `json:"title" db:"title"` // lower and trim
	Body      string    `json:"body" db:"body"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	ValidTill time.Time `json:"validTill" db:"valid_till"`
	Status    string    `json:"status" db:"status"`
}

type Session struct {
	SessionID  string    `json:"-" db:"session_id"`
	UserID     string    `json:"user_id" db:"user_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	ArchivedAt time.Time `json:"archived_at" db:"archived_at"`
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=8,lte=30"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CreateTodoRequest struct {
	Title     string    `json:"title" validate:"required,min=3,max=100"`
	Body      string    `json:"body" validate:"required"`
	ValidTill time.Time `json:"validTill" validate:"required"`
}

type UpdateTodoRequest struct {
	Title     *string    `json:"title" validate:"omitempty,min=3,max=100"`
	Body      *string    `json:"body" validate:"omitempty"`
	Status    *string    `json:"status"`
	ValidTill *time.Time `json:"valid_till"`
}

type UpcomingTodosRequest struct {
	Title     *string    `json:"title"`
	Body      *string    `json:"body"`
	Status    *string    `json:"status"`
	ValidTill *time.Time `json:"valid_till"`
}

type RequestContext struct {
	UserID    string `json:"userID"`
	SessionID string `json:"sessionID"`
}

type Error struct {
	Error      string `json:"error"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}
