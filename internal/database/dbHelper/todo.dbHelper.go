package dbhelper

import (
	"time"

	"github.com/Saikatdeb12/TodoApp2/internal/database"
	"github.com/Saikatdeb12/TodoApp2/internal/models"
)

func ArchiveSession(sessionID string) (int64, error) {
	SQL := `
		UPDATE sessions
		SET archived_at = NOW()
		WHERE id = $1 AND archived_at IS NULL
	`
	res, err := database.DB.Exec(SQL, sessionID)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func CreateTodoSQL(user_id string, title, body string, valid_till time.Time) (models.Todo, error) {
	SQL := `
		INSERT INTO todos (user_id, title, body, valid_till)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, title, body, created_at, valid_till, status
	`
	var todo models.Todo
	err := database.DB.Get(&todo, SQL, user_id, title, body, valid_till)
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
	err := database.DB.Select(&todos, SQL, args...)
	return todos, err
}

func GetAllTodosByFilter(user_id string, status string) ([]models.Todo, error) {
	SQL := `
		SELECT id, title, body, created_at, valid_till, status
		FROM todos
		WHERE user_id=$1 AND status=$2 AND archived_at IS NULL
		ORDER BY created_at DESC
	`

	todos := []models.Todo{}
	rows, err := database.DB.Query(SQL, user_id, status)
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

func GetTodoByID(userID, todoID string) (*models.Todo, error) {
	SQL := `
		SELECT id, user_id, title, body, created_at, valid_till, status
		FROM todos
		WHERE id=$1 AND user_id=$2 AND archived_at is NULL;
	`
	todo := models.Todo{}

	// err := DB.QueryRow(SQL, todo_ID, user_ID).Scan(&todo.TodoID, &todo.UserID, &todo.Title, &todo.Body, &todo.CreatedAt, &todo.ValidTill, &todo.Complete)
	// not to use QueryRow
	err := database.DB.Get(&todo, SQL, todoID, userID)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func UpdateTodoById(name, description string, status string, expiringAt time.Time, todoID, userID string) error {
	SQL := `UPDATE todos 
			SET name=$1,description=$2,status=$3,expiring_at=$4
			WHERE id=$5 
			and user_id=$6;`
	_, err := database.DB.Exec(
		SQL, name, description, status, expiringAt, todoID, userID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteTodoByID(userID, todoID string) (int64, error) {
	query := `
		UPDATE todos
		SET archived_at=NOW()
		WHERE user_id=$1
		AND todo_id=$2
		AND archived_at is NULL
	`
	res, err := database.DB.Exec(query, todoID, userID)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
