package dbhelper

import (
	"database/sql"

	"github.com/Saikatdeb12/TodoApp2/internal/database"
	"github.com/Saikatdeb12/TodoApp2/internal/models"
)

func CheckUserExistsByEmail(email string) (bool, error) {
	var id int
	SQL := `
		SELECT id FROM users
		WHERE email = $1
		AND archived_at IS NULL
	`
	err := database.DB.Get(&id, SQL, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func ValidateUserSession(sessionID string) (string, error) {
	SQL := `
		SELECT user_id FROM session 
		WHERE id=$1 
		AND archived_at is NULL
	`

	// var currentSession models.Session
	var userID string

	err := database.DB.Get(&userID, SQL, sessionID)
	if err != nil {
		return "", err
	}

	// return currentSession.UserID, nil
	return userID, nil
}

func CreateUser(name, email, password string) (string, error) {
	SQL := `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var userID string
	err := database.DB.Get(&userID, SQL, name, email, password)
	return userID, err
}

func GetUserAuthByEmail(email string) (models.User, error) {
	var user models.User
	SQL := `
		SELECT id, password, name, email, created_at, archived_at 
		FROM users
		WHERE email=$1 AND archived_at IS NULL
	`

	err := database.DB.Get(&user, SQL, email)
	return user, err
}

func CreateSession(userID string) (string, error) {
	SQL := `
		INSERT INTO sessions (user_id) 
		VALUES($1)
		RETURNING id;
	`
	var sessionID string
	err := database.DB.Get(&sessionID, SQL, userID)
	return sessionID, err
}

func ArchiveUser(userID string) (int64, error) {
	SQL := `
		UPDATE users
		SET archived_at = NOW()
		WHERE id = $1 AND archived_at IS NULL
	`

	res, err := database.DB.Exec(SQL, userID)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
