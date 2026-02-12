package dbhelper

import (
	"database/sql"

	"github.com/Saikatdeb12/TodoApp2/internal/database"
	"github.com/Saikatdeb12/TodoApp2/internal/models"
	"github.com/jmoiron/sqlx"
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

func ValidateUserSession(sessionID string) error {
	// SQL := `
	// 	SELECT user_id FROM session
	// 	WHERE id=$1
	// 	AND archived_at is NULL
	// `
	SQL := `
		UPDATE sessions
		SET archived_at=NOW()
		WHERE id=$1
		AND archived_at IS NULL;
	`
	err := database.DB.Get(SQL, sessionID)
	return err
}

func CreateUserTx(trx *sqlx.Tx, name, email, password string) (string, error) {
	SQL := `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var userID string
	err := trx.Get(&userID, SQL, name, email, password)
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

func CreateSessionTx(trx *sqlx.Tx, userID string) (string, error) {
	SQL := `
		INSERT INTO sessions (user_id) 
		VALUES($1)
		RETURNING id;
	`
	var sessionID string
	err := trx.Get(&sessionID, SQL, userID)
	return sessionID, err
}

func ArchiveUser(userID string) error {
	SQL := `
		UPDATE users
		SET archived_at = NOW()
		WHERE id = $1 AND archived_at IS NULL
	`

	_, err := database.DB.Exec(SQL, userID)
	return err
}

func ArchiveSession(userID string) error {
	SQL := `
		UPDATE sessions
		SET archived_at = NOW()
		WHERE id = $1 AND archived_at IS NULL
	`

	_, err := database.DB.Exec(SQL, userID)
	return err
}

func ArchiveUserTx(trx *sqlx.Tx, userID string) error {
	SQL := `
		UPDATE users
		SET archived_at = NOW()
		WHERE id = $1 AND archived_at IS NULL
	`
	_, err := trx.Exec(SQL, userID)
	return err
}

func ArchiveSessionTx(trx *sqlx.Tx, userID string) error {
	SQL := `
		UPDATE sessions
		SET archived_at = NOW()
		WHERE id = $1 AND archived_at IS NULL
	`

	_, err := trx.Exec(SQL, userID)
	return err
}
