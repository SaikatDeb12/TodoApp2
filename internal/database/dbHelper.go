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

func CreateTodoSQL(user_id uuid.UUID, title, body string, valid_till time.Time) (models.Todo, error) {
	query := `
		INSERT INTO todos (user_id, title, body, valid_till)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, title, body, created_at, valid_till, complete
	`
	var todo models.Todo
	err := DB.QueryRow(
		query, user_id, title, body, valid_till,
	).Scan(&todo.TodoID, &todo.UserID, &todo.Title, &todo.Body, &todo.CreatedAt, &todo.ValidTill, &todo.Complete)

	return todo, err
}

func GetAllTodos(user_id uuid.UUID) ([]models.Todo, error) {
	SQL := `
		SELECT id, title, body, created_at, valid_till
		FROM todos
		WHERE user_id=$1
	`
	todos := []models.Todo{}
	rows, err := DB.Query(SQL, user_id)
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

func GetAllTodosByFilter(user_id uuid.UUID, complete bool) ([]models.Todo, error) {
	SQL := `
		SELECT id, title, body, created_at, valid_till
		FROM todos
		WHERE user_id=$1 AND complete=$2
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
		if err := rows.Scan(&todo.TodoID, &todo.Title, &todo.Body, &todo.CreatedAt, &todo.ValidTill); err != nil {
			return todos, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}
