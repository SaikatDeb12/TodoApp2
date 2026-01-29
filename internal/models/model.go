package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID uuid.UUID `json:"id" db:"user_id"`
	Name string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
	Password string `json:"Password" db:"password"`
	
}

type Todo struct {
	TodoID uuid.UUID `json:"id" db:"todo_id"`
	UserID uuid.UUID `json:"-" db:"user_id"`
	Title string `json:"title" db:"title"`
	Body string `json:"body" db:"body"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	ValidTill time.Time `json:"validTill" db:"valid_till"`
	Completed  bool  `json:"complete" db:"complete"`
}

type Session struct{
	SessionID uuid.UUID `db:"session_id"`
	UserID uuid.UUID `db:"session_id"`
	CreatedAt time.Time `db:"session_id"`
	ExpiresAt time.Time `db:"session_id"`
}
