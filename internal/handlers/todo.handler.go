package handlers

import (
	"net/http"
	"time"

	dbhelper "github.com/Saikatdeb12/TodoApp2/internal/database/dbHelper"
	"github.com/Saikatdeb12/TodoApp2/internal/middleware"
	middlewares "github.com/Saikatdeb12/TodoApp2/internal/middleware"
	"github.com/Saikatdeb12/TodoApp2/internal/models"
	"github.com/Saikatdeb12/TodoApp2/internal/utils"
	"github.com/go-chi/chi/v5"
)

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTodoRequest
	if err := utils.ParseBody(r.Body, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "invalid payload")
		return
	}

	userContext, ok := middlewares.UserContext(r)
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
		return
	}

	userID := userContext.UserID
	var todo models.Todo
	todo, err := dbhelper.CreateTodoSQL(userID, req.Title, req.Body, req.ValidTill)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "internal server error")
		return
	}

	// can do without sending created todo in response
	todo.UserID = userID
	todo.Title = req.Title
	todo.Body = req.Body
	todo.ValidTill = req.ValidTill

	utils.RespondJSON(w, http.StatusCreated, todo)
}

func GetTodos(w http.ResponseWriter, r *http.Request) {
	completeStr := r.URL.Query().Get("status")
	expiringAtStr := r.URL.Query().Get("expiringAt")
	search := r.URL.Query().Get("search")

	userCtx, _ := middleware.UserContext(r)
	userID := userCtx.UserID
	if expiringAtStr != "" {
		d, err := time.Parse("2006-01-02", expiringAtStr)
		if err != nil {
			utils.RespondError(w, http.StatusBadRequest, err, "invalid date")
			return
		}
		if d.Before(time.Now()) {
			expiringAtStr = ""
		}
	}
	todos, err := dbhelper.GetTodos(userID, search, expiringAtStr, completeStr)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Failed to fetch todos")
		return
	}

	utils.RespondJSON(w, http.StatusOK, struct {
		Todos []models.Todo `json:"todos"`
	}{
		Todos: todos,
	})
}

func GetTodoByID(w http.ResponseWriter, r *http.Request) {
	userContext, ok := middlewares.UserContext(r)
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
		return
	}
	userID := userContext.UserID

	todoID := chi.URLParam(r, "id")

	todo, err := dbhelper.GetTodoByID(userID, todoID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err, "todo not found")
		return
	}

	utils.RespondJSON(w, http.StatusOK, todo)
}

func UpdateTodoByID(w http.ResponseWriter, r *http.Request) {
	todoID := chi.URLParam(r, "id")

	var todo models.UpdateTodoRequest
	if err := utils.ParseBody(r.Body, &todo); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "invalid request body")
		return
	}

	if err := utils.ValidateStruct(todo); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "validation failed")
		return
	}

	userContext, ok := middlewares.UserContext(r)
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
		return
	}
	userID := userContext.UserID

	err := dbhelper.UpdateTodoById(*todo.Title, *todo.Body, *todo.Status, *todo.ValidTill, todoID, userID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "invalid payload")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "todo updated successfully",
	})
}

func DeleteTodoByID(w http.ResponseWriter, r *http.Request) {
	userContext, ok := middlewares.UserContext(r)
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
		return
	}
	userID := userContext.UserID

	todoID := chi.URLParam(r, "id")

	affected, err := dbhelper.DeleteTodoByID(userID, todoID)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "delete operation failed")
		return
	}

	if affected == 0 {
		utils.RespondError(w, http.StatusNotFound, err, "todo not found")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Todo deleted",
	})
}
