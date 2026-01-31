package database

import (
	"time"

	"github.com/Saikatdeb12/TodoApp2/internal/models"
	"github.com/google/uuid"
)

func CheckUserExistsByEmail(email string) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS(
			SELECT 1 FROM users 
			WHERE email=$1 AND archived_at IS NULL
		)
	`
	err := DB.QueryRow(query, email).Scan(&exists)
	return exists, err
}

func CreateUser(name, email, password string) error {
	query := `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
	`
	_, err := DB.Exec(query, name, email, password)
	return err
}

func GetUserAuthByEmail(email string) (models.User, error) {
	var user models.User
	query := `SELECT id, password FROM users WHERE email=$1 AND archived_at IS NULL`

	err := DB.QueryRow(query, email).Scan(&user.UserID, &user.Password)
	return user, err
}

func CreateSession(sessionID, userID uuid.UUID, expires time.Time) error {
	query := `
		INSERT INTO sessions (id, user_id, created_at, expires_at)
		VALUES($1, $2, now(), $3)
	`
	_, err := DB.Exec(query, sessionID, userID, expires)
	return err
}

func DeleteSession(sessionID uuid.UUID) (int64, error) {
	query := `DELETE FROM sessions WHERE id=$1`
	res, err := DB.Exec(query, sessionID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
