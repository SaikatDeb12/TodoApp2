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
	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusUnauthorized)
		return
	}

	todoID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid todo id", http.StatusBadRequest)
		return
	}

	req, err := database.ParseUpdateTodoRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	affected, err := database.UpdateTodoByID(userID, todoID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if affected == 0 {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
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
		http.Error(w, "Invalid todo id", http.StatusBadGateway)
		return
	}

	affected, err := database.DeleteTodoByID(userID, todoID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if affected == 0 {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Todo deleted",
	})
}

func UpcomingTodosByDate(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusUnauthorized)
		return
	}
	// exec,
	days := 0
	if dayParam := r.URL.Query().Get("days"); dayParam != "" {
		days, err = strconv.Atoi(dayParam)
		if err != nil {
			http.Error(w, "Invalid days parameter", http.StatusBadRequest)
			return
		}
	}

	todos, err := database.GetUpcomingTodos(userID, days)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.RespondJSON(w, http.StatusOK, todos)
}
