package database

import (
	"database/sql"
	"time"

	"github.com/Saikatdeb12/TodoApp2/internal/models"
	"github.com/google/uuid"
)

func CheckUserExistsByEmail(email string) (bool, error) {
	var id int
	SQL := `
		SELECT * FROM users
		WHERE email = $1
		AND archived_at IS NULL
		LIMIT 1
	`
	err := DB.Get(&id, SQL, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func ValidateUserSession(sessionID uuid.UUID) (uuid.UUID, error) {
	SQL := `
		SELECT user_id FROM session 
		WHERE id=$1 
		AND archived_at is NULL
		AND expires_at > NOW()
	`

	var currentSession models.Session

	err := DB.Get(&currentSession, SQL, sessionID)
	if err != nil {
		return uuid.Nil, err
	}

	return currentSession.UserID, nil
}

func CreateUser(name, email, password string) (uuid.UUID, error) {
	SQL := `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var userID uuid.UUID
	err := DB.Get(&userID, SQL, name, email, password)
	return userID, err
}

func GetUserAuthByEmail(email string) (models.User, error) {
	var user models.User
	SQL := `SELECT id, password FROM users WHERE email=$1 AND archived_at IS NULL`

	err := DB.Get(&user, SQL, email)
	return user, err
}

func CreateSession(userID uuid.UUID) (string, error) {
	SQL := `
		INSERT INTO sessions (id, user_id, created_at, expires_at)
		VALUES($1, $2, now(), $3)
		RETURNING id;
	`
	var sessionID string
	err := DB.Get(&sessionID, SQL, userID)
	if err != nil {
		return "", err
	}
	return sessionID, nil
}

func ArchiveUser(userID uuid.UUID) (int64, error) {
	SQL := `
		UPDATE users
		SET archived_at = NOW()
		WHERE id = $1 AND archived_at IS NULL
	`

	res, err := DB.Exec(SQL, userID)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func ArchiveSession(sessionID uuid.UUID) (int64, error) {
	SQL := `
		UPDATE sessions
		SET archived_at = NOW()
		WHERE id = $1 AND archived_at IS NULL
	`
	res, err := DB.Exec(SQL, sessionID)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func CreateTodoSQL(user_id uuid.UUID, title, body string, valid_till time.Time) (models.Todo, error) {
	SQL := `
		INSERT INTO todos (user_id, title, body, valid_till)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, title, body, created_at, valid_till, status
	`
	var todo models.Todo
	err := DB.Get(&todo, SQL, user_id, title, body, valid_till)
	return todo, err
}

func GetAllTodos(userID, status string, selectedDate *time.Time) ([]models.Todo, error) {
	SQL := `SELECT todo_id, title, description, status, valid_till, created_at
             FROM todos
             WHERE user_id = $1
               AND archived_at IS NULL
               AND (
                   $2 = ''
                   OR status = $2::todo_status
                  )
               AND (
                   $3::timestamp IS NULL OR valid_till <= $3::timestamp
                  )
               order by created_at desc;`

	args := []interface{}{userID, status, selectedDate}

	todos := make([]models.Todo, 0)
	err := DB.Select(&todos, SQL, args...)
	return todos, err
}

func GetAllTodosByFilter(user_id uuid.UUID, status string) ([]models.Todo, error) {
	SQL := `
		SELECT id, title, body, created_at, valid_till, status
		FROM todos
		WHERE user_id=$1 AND status=$2 AND archived_at IS NULL
		ORDER BY created_at DESC
	`

	todos := []models.Todo{}
	rows, err := DB.Query(SQL, user_id, status)
	if err != nil {
		return todos, err
	}
	defer rows.Close()

	for rows.Next() {
		var todo models.Todo
		if err := rows.Scan(&todo.TodoID, &todo.Title, &todo.Body, &todo.CreatedAt, &todo.ValidTill, &todo.Status); err != nil {
			return todos, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

func GetTodoByID(userID, todoID uuid.UUID) (*models.Todo, error) {
	SQL := `
		SELECT id, user_id, title, body, created_at, valid_till, status
		FROM todos
		WHERE id=$1 AND user_id=$2 AND archived_at is NULL;
	`
	todo := models.Todo{}

	// err := DB.QueryRow(SQL, todo_ID, user_ID).Scan(&todo.TodoID, &todo.UserID, &todo.Title, &todo.Body, &todo.CreatedAt, &todo.ValidTill, &todo.Complete)
	err := DB.Get(&todo, SQL, todoID, userID)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

type UpdateTodoRequest struct {
	Title     *string    `json:"title"`
	Body      *string    `json:"body"`
	Status    *string    `json:"status"`
	ValidTill *time.Time `json:"valid_till"`
}

func UpdateTodoById(name, description string, status string, expiringAt time.Time, todoID, userID uuid.UUID) error {
	SQL := `UPDATE todos 
			SET name=$1,description=$2,status=$3,expiring_at=$4
			WHERE id=$5 
			and user_id=$6;`
	_, err := DB.Exec(
		SQL, name, description, status, expiringAt, todoID, userID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteTodoByID(userID, todoID uuid.UUID) (int64, error) {
	query := `
		UPDATE todos
		SET archived_at=NOW()
		WHERE user_id=$1
		AND todo_id=$2
		AND archived_at is NULL
	`
	res, err := DB.Exec(query, todoID, userID)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

// func GetUpcomingTodos(userID uuid.UUID, days int) ([]models.Todo, error) {
// 	query := `
// 		SELECT id, title, body, created_at, status, valid_till
// 		FROM todos
// 		WHERE user_id = $1
// 		  AND status = 'incomplete'
// 		  AND valid_till IS NOT NULL
// 		  AND valid_till BETWEEN CURRENT_DATE
// 		  AND CURRENT_DATE + ($2 || ' days')::INTERVAL
// 		  AND archived_at IS NULL
// 		ORDER BY valid_till;
// 	`
//
// 	rows, err := DB.Query(query, userID, days)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
//
// 	var todos []models.Todo
// 	for rows.Next() {
// 		var todo models.Todo
// 		if err := rows.Scan(
// 			&todo.TodoID,
// 			&todo.Title,
// 			&todo.Body,
// 			&todo.CreatedAt,
// 			&todo.Status,
// 			&todo.ValidTill,
// 		); err != nil {
// 			return nil, err
// 		}
// 		todos = append(todos, todo)
// 	}
//
// 	return todos, nil
// }
