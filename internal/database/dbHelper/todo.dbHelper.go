package dbhelper

import (
	"time"

	"github.com/Saikatdeb12/TodoApp2/internal/database"
	"github.com/Saikatdeb12/TodoApp2/internal/models"
	"github.com/jmoiron/sqlx"
)

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

func GetTodos(userID, title, date, status string, limit, offset int) ([]models.Todo, error) {
	SQL := `
          SELECT id,
                 user_id,
                 title,
                 body,
                 status,
                 valid_till,
                created_at
          FROM todos
          WHERE user_id =$1
          AND ( $2 = '' or status=$2::todo_status)
          AND ( $3='' or valid_till<=$3::TIMESTAMPTZ)
          AND ( $4::TEXT IS NULL OR title LIKE'%'||$4||'%')
		  AND archived_at IS NULL
          ORDER BY valid_till
		  LIMIT $5 OFFSET $6
          `
	todos := make([]models.Todo, 0)

	err := database.DB.Select(&todos, SQL, userID, status, date, title, limit, offset)
	return todos, err
}

func GetAllTodosByFilter(user_id string, status string) ([]models.Todo, error) {
	SQL := `
		SELECT id, title, body, created_at, valid_till, status
		FROM todos
		WHERE user_id=$1 AND status=$2 AND archived_at IS NULL
		ORDER BY created_at DESC
	`

	todos := make([]models.Todo, 0)
	err := database.DB.Select(&todos, SQL, user_id, status)
	if err != nil {
		return nil, err
	}

	return todos, nil
}

func GetTodoByID(userID, todoID string) (*models.Todo, error) {
	SQL := `
		SELECT id, user_id, title, body, created_at, valid_till, status
		FROM todos
		WHERE id=$1 AND user_id=$2 AND archived_at IS NOT NULL;
	`
	// err := DB.QueryRow(SQL, todo_ID, user_ID).Scan(&todo.TodoID, &todo.UserID, &todo.Title, &todo.Body, &todo.CreatedAt, &todo.ValidTill, &todo.Complete)
	// not to use QueryRow

	var todo models.Todo
	err := database.DB.Get(&todo, SQL, todoID, userID)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func UpdateTodoById(title, body, status, valid_till *string, todoID, userID string) error {
	SQL := `UPDATE todos 
			SET title=COALESCE($1,title),
			body=COALESCE($2,body),
			status=COALESCE($3,status),
			valid_till=COALESCE($4::timestamptz,valid_till)
			WHERE id=$5 
			and user_id=$6;`
	_, err := database.DB.Exec(
		SQL, title, body, status, valid_till, todoID, userID)
	return err
}

func DeleteTodoByID(userID, todoID string) error {
	SQL := `
		UPDATE todos
		SET archived_at=NOW()
		WHERE user_id=$1
		AND id=$2
		AND archived_at is NULL
	`
	_, err := database.DB.Exec(SQL, userID, todoID)

	return err
}

func DeleteAllTodos(userID string) error {
	SQL := `
		UPDATE todos
		SET archived_at=NOW()
		WHERE user_id=$1
		AND archived_at is NULL
	`
	_, err := database.DB.Exec(SQL, userID)
	return err
}

func DeleteAllTodosTx(trx *sqlx.Tx, userID string) error {
	SQL := `
		UPDATE todos
		SET archived_at=NOW()
		WHERE user_id=$1
		AND archived_at is NULL
	`
	_, err := trx.Exec(SQL, userID)
	return err
}
