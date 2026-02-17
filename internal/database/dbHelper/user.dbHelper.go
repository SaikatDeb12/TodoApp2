package dbhelper

import (
	"github.com/Saikatdeb12/TodoApp2/internal/database"
	"github.com/Saikatdeb12/TodoApp2/internal/models"
	"github.com/jmoiron/sqlx"
)

func CheckUserExistsByEmail(email string) error {
	SQL := `
		SELECT COUNT(*) 
		FROM users
		WHERE email=TRIM(LOWER($1)) AND archived_at IS NULL
	`
	err := database.DB.Get(SQL, email)
	return err
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
	_, err := database.DB.Exec(SQL, sessionID)
	return err
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

func ArchiveSession(userID string) error {
	SQL := `
		UPDATE sessions
		SET archived_at = NOW()
		WHERE id = $1 AND archived_at IS NULL
	`
	_, err := database.DB.Exec(SQL, userID)
	return err
}

//func ArchiveUser(userID string) error {
//	SQL := `
//		UPDATE users
//		SET archived_at = NOW()
//		WHERE id = $1 AND archived_at IS NULL
//	`
//	_, err := database.DB.Exec(SQL, userID)
//	return err
//}

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
