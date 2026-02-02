package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Saikatdeb12/TodoApp2/internal/database"
	"github.com/Saikatdeb12/TodoApp2/internal/models"
	"github.com/Saikatdeb12/TodoApp2/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var todo models.Todo
	todo, err = database.CreateTodoSQL(userID, req.Title, req.Body, req.ValidTill)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	todo.UserID = userID
	todo.Title = req.Title
	todo.Body = req.Body
	todo.ValidTill = req.ValidTill

	utils.RespondJSON(w, http.StatusCreated, todo)
}

func GetTodos(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	query := r.URL.Query().Get("complete")

	var todos []models.Todo
	if query == "" {
		todos, err = database.GetAllTodos(userID)
	} else {
		complete, err := strconv.ParseBool(query)
		if err != nil {
			http.Error(w, "Invalid query", http.StatusBadRequest)
			return
		}
		todos, err = database.GetAllTodosByFilter(userID, complete)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.RespondJSON(w, http.StatusOK, todos)
}

func GetTodoByID(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	todoID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusInternalServerError)
		return
	}

	todo, err := database.GetTodoByID(userID, todoID)
	if err != nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	utils.RespondJSON(w, http.StatusOK, todo)
}

func UpdateTodoByID(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUserID(r.Context())
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusUnauthorized)
		return
	}

	todoId, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid todo id", http.StatusBadRequest)
		return
	}

	var req models.UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// again to check the number of affect rows
	aff, _ := res.RowsAffected()
	if aff == 0 {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
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
