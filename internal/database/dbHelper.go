package database

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Saikatdeb12/TodoApp2/internal/models"
	"github.com/google/uuid"
)

func CheckUserExistsByEmail(email string) (bool, error) {
	var exists bool
	SQL := `
		SELECT EXISTS(
			SELECT 1 FROM users 
			WHERE email=$1 AND archived_at IS NULL
		)
	`
	err := DB.Get(&exists, SQL, email)
	return exists, err
}

func ValidateUserSession(sessionID string) (uuid.UUID, error) {
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

func CreateSession(sessionID, userID uuid.UUID, expires time.Time) error {
	SQL := `
		INSERT INTO sessions (id, user_id, created_at, expires_at)
		VALUES($1, $2, now(), $3)
	`
	_, err := DB.Exec(SQL, sessionID, userID, expires)
	return err
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

// func DeleteSession(sessionID uuid.UUID) (int64, error) {
// 	SQL := `DELETE FROM sessions WHERE id=$1`
// 	res, err := DB.Exec(SQL, sessionID)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return res.RowsAffected()
// }

func CreateTodoSQL(user_id uuid.UUID, title, body string, valid_till time.Time) (models.Todo, error) {
	SQL := `
		INSERT INTO todos (user_id, title, body, valid_till)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, title, body, created_at, valid_till, complete
	`
	var todo models.Todo
	err := DB.Get(&todo, SQL, user_id, title, body, valid_till)
	return todo, err
}

func GetAllTodos(user_id uuid.UUID) ([]models.Todo, error) {
	SQL := `
		SELECT id,user_id, title, body, created_at, valid_till, complete
		FROM todos
		WHERE user_id=$1 AND archived_at is NULL
	`
	todos := []models.Todo{}
	rows, err := DB.Query(SQL, user_id)
	if err != nil {
		return todos, err
	}

	defer rows.Close()

	for rows.Next() {
		var todo models.Todo
		if err := rows.Scan(&todo.TodoID, &todo.UserID, &todo.Title, &todo.Body, &todo.CreatedAt, &todo.ValidTill, &todo.Complete); err != nil {
			return todos, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

func GetAllTodosByFilter(user_id uuid.UUID, complete bool) ([]models.Todo, error) {
	SQL := `
		SELECT id, title, body, created_at, valid_till, complete
		FROM todos
		WHERE user_id=$1 AND complete=$2 AND archived_at IS NULL
		ORDER BY created_at DESC
	`

	todos := []models.Todo{}
	rows, err := DB.Query(SQL, user_id, complete)
	if err != nil {
		return todos, err
	}
	defer rows.Close()

	for rows.Next() {
		var todo models.Todo
		if err := rows.Scan(&todo.TodoID, &todo.Title, &todo.Body, &todo.CreatedAt, &todo.ValidTill, &todo.Complete); err != nil {
			return todos, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

func GetTodoByID(userID, todoID uuid.UUID) (*models.Todo, error) {
	SQL := `
		SELECT id, user_id, title, body, created_at, valid_till, complete
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
	Complete  *bool      `json:"complete"`
	ValidTill *time.Time `json:"valid_till"`
}

func ParseUpdateTodoRequest(r *http.Request) (*UpdateTodoRequest, error) {
	var req UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.New("invalid payload")
	}
	return &req, nil
}

func UpdateTodoByID(
	userID uuid.UUID,
	todoID uuid.UUID,
	req *UpdateTodoRequest,
) (int64, error) {
	query := `
		UPDATE todos
		SET
			title      = COALESCE($1, title),
			body       = COALESCE($2, body),
			complete   = COALESCE($3, complete),
			valid_till = COALESCE($4, valid_till)
		WHERE id = $5 AND user_id = $6
	`

	res, err := DB.Exec(
		query,
		req.Title,
		req.Body,
		req.Complete,
		req.ValidTill,
		todoID,
		userID,
	)
	if err != nil {
		return 0, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affected, nil
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

func GetUpcomingTodos(userID uuid.UUID, days int) ([]models.Todo, error) {
	query := `
		SELECT id, title, body, created_at, complete, valid_till
		FROM todos
		WHERE user_id = $1
		  AND complete = false
		  AND valid_till IS NOT NULL
		  AND valid_till BETWEEN CURRENT_DATE
		  AND CURRENT_DATE + ($2 || ' days')::INTERVAL
		  AND archived_at IS NULL
		ORDER BY valid_till;
	`

	rows, err := DB.Query(query, userID, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		if err := rows.Scan(
			&todo.TodoID,
			&todo.Title,
			&todo.Body,
			&todo.CreatedAt,
			&todo.Complete,
			&todo.ValidTill,
		); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

// TODO defining struct for all the function using POST method
