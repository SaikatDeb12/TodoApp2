package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Saikatdeb12/TodoApp2/internal/database"
	"github.com/Saikatdeb12/TodoApp2/internal/models"
	"github.com/Saikatdeb12/TodoApp2/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CreateTodoRequest struct {
	Title string `json:"title"`
	Body string `json:"body"`
	ValidTill time.Time `json:"validTill"`
}

func CreateTodo(w http.ResponseWriter, r *http.Request){
	var req CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	userId, err := utils.GetUserID(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	query := `
		INSERT INTO todos (user_id, title, body, valid_till)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	var todo models.Todo
	err = database.DB.QueryRow(
		query, userId, req.Title, req.Body, req.ValidTill,
	).Scan(&todo.TodoID, &todo.CreatedAt)

	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	todo.Title=req.Title
	todo.Body=req.Body
	todo.ValidTill=req.ValidTill
		
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)

}

func GetTodos(w http.ResponseWriter, r *http.Request){
	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return 
	}

	query := `
		SELECT id, title, body, created_at, complete, valid_till
		FROM todos
		WHERE user_id=$1
		ORDER BY created_at DESC
	`

	rows, err := database.DB.Query(query, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	todos := []models.Todo{}
	for rows.Next(){
		var todo models.Todo
		if err := rows.Scan(&todo.TodoID, &todo.Title, &todo.Body, &todo.CreatedAt, &todo.Complete, &todo.ValidTill); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		todos = append(todos, todo)
	}

	json.NewEncoder(w).Encode(todos)
}

func GetTodoByID (w http.ResponseWriter, r *http.Request){
	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	todoID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid todo id", http.StatusBadRequest)
		return
	}

	query := `
		SELECT id, title, body, created_at, complete, valid_till
		FROM todos
		WHERE id=$1 AND user_id=$2
	`

	var todo models.Todo

	err = database.DB.QueryRow(query, todoID, userID).Scan(
		&todo.TodoID,
		&todo.Title,
		&todo.Body,
		&todo.CreatedAt,
		&todo.Complete,
		&todo.ValidTill,
	)

	if err == sql.ErrNoRows {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(todo)
}

type UpdateTodoRequest struct{
	Title *string `json:"title"`
	Body *string `json:"body"`
	Complete *bool `json:"complete"`
	ValidTill *time.Time `json:"valid_till"`
}

func UpdateTodoByID(w http.ResponseWriter, r *http.Request){
	userId, err := utils.GetUserID(r.Context())
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusUnauthorized)
		return
	}

	todoId, err := uuid.Parse(chi.URLParam(r, "id"))
	if err !=nil {
		http.Error(w, "Invalid todo id", http.StatusBadRequest)
		return
	}

	var req UpdateTodoRequest
	if err:= json.NewDecoder(r.Body).Decode(&req); err!=nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	query := `
		UPDATE todos
		SET
			title=coalesce($1, title),
			body=coalesce($2, body),
			complete=coalesce($3, complete),
			valid_till=coalesce($4, valid_till)
		WHERE id=$5 AND user_id=$6
	`

	res, err := database.DB.Exec(query, req.Title, req.Body, req.Complete, req.ValidTill, todoId, userId)
	if err!= nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//again to check the number of affect rows
	aff, _ := res.RowsAffected()
	if aff == 0{
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string] string {
		"msg": "Todo updated successfully",
	})

}

func DeleteTodoByID(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	todoID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid todo id", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM todos WHERE id=$1 AND user_id=$2`

	res, err := database.DB.Exec(query, todoID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	aff, _ := res.RowsAffected()
	if aff == 0 {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"msg": "Todo deleted"})
}

func CompletedTodos(w http.ResponseWriter, r *http.Request){
	userId, err := utils.GetUserID(r.Context())
	if err != nil {
		http.Error(w, "Invalid user", http.StatusUnauthorized)
		return 
	}

	query := `
		SELECT id, title, body, created_at, complete, valid_till
		FROM todos
		WHERE user_id=$1 and complete=true
		ORDER BY created_at DESC
	`

	rows, err := database.DB.Query(query, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}

	defer rows.Close()


	todos := []models.Todo{}
	for rows.Next(){
		var todo models.Todo
		if err := rows.Scan(&todo.TodoID, &todo.Title, &todo.Body, &todo.CreatedAt, &todo.Complete, &todo.ValidTill); err!=nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		todos=append(todos, todo)
	}
	
	json.NewEncoder(w).Encode(todos)
}
func InCompleteTodos(w http.ResponseWriter, r *http.Request){
	userId, err := utils.GetUserID(r.Context())
	if err != nil {
		http.Error(w, "Invalid user", http.StatusUnauthorized)
		return 
	}

	query := `
		SELECT id, title, body, created_at, complete, valid_till
		FROM todos
		WHERE user_id=$1 and complete=false
		ORDER BY created_at DESC
	`

	rows, err := database.DB.Query(query, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}

	defer rows.Close()

	todos := []models.Todo{}
	for rows.Next(){
		var todo models.Todo
		if err := rows.Scan(&todo.TodoID, &todo.Title, &todo.Body, &todo.CreatedAt, &todo.Complete, &todo.ValidTill); err!=nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		todos=append(todos, todo)
	}
	
	json.NewEncoder(w).Encode(todos)

}

type UpcomingTodosRequest struct{
	Title *string `json:"title"`
	Body *string `json:"body"`
	Complete *bool `json:"complete"`
	ValidTill *time.Time `json:"valid_till"`
}

func UpcomingTodosByDate(w http.ResponseWriter, r *http.Request){
	userId, err := utils.GetUserID(r.Context())
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusUnauthorized)
		return
	}

	dayParam:=r.URL.Query().Get("days")
	var days int
	if(dayParam==""){
		days=0
	} else{
		days, err = strconv.Atoi(r.URL.Query().Get("days"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return;
		}
	}

	query := `
		SELECT id, title, body, created_at, complete, valid_till
		FROM todos
		WHERE user_id = $1
		  AND complete = false
		  AND valid_till IS NOT NULL
		  AND valid_till BETWEEN CURRENT_DATE
		  AND CURRENT_DATE + ($2 || ' days')::INTERVAL
		ORDER BY valid_till;
	`
	rows, err := database.DB.Query(query, userId, days)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var todos []models.Todo
	for rows.Next(){
		var todo models.Todo
		if err := rows.Scan(
			&todo.TodoID, &todo.Title, &todo.Body, &todo.CreatedAt, &todo.Complete, &todo.ValidTill,
		); err != nil {
			http.Error(w, "No todos found", http.StatusNotFound)
			return
		}

		todos = append(todos, todo)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todos)
}



