package handlers

import (
	"net/http"

	"github.com/Saikatdeb12/TodoApp2/internal/database"
	middlewares "github.com/Saikatdeb12/TodoApp2/internal/middleware"
	"github.com/Saikatdeb12/TodoApp2/internal/models"
	"github.com/Saikatdeb12/TodoApp2/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTodoRequest
	if err := utils.ParseBody(r.Body, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "invalid payload")
		return
	}

	userContext := middlewares.UserContext(r)
	userID := userContext.UserID
	var todo models.Todo
	todo, err := database.CreateTodoSQL(userID, req.Title, req.Body, req.ValidTill)
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
	userContext := middlewares.UserContext(r)
	userID := userContext.UserID

	status := r.URL.Query().Get("status")
	expiryDate := r.URL.Query().Get("date")

	date, err := utils.ParseDate(expiryDate)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to parse date")
		return
	}

	todos, getErr := database.GetAllTodos(userID.String(), status, date)

	if getErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, getErr, "failed to fetch todos")
		return
	}

	utils.RespondJSON(w, http.StatusOK, todos)
}

func GetTodoByID(w http.ResponseWriter, r *http.Request) {
	userContext := middlewares.UserContext(r)
	userID := userContext.UserID

	todoID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "valid todo ID")
		return
	}

	todo, err := database.GetTodoByID(userID, todoID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err, "todo not found")
		return
	}

	utils.RespondJSON(w, http.StatusOK, todo)
}

func UpdateTodoByID(w http.ResponseWriter, r *http.Request) {
	todoID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "invalid todo id")
		return
	}

	var todo models.UpdateTodoRequest
	if err := utils.ParseBody(r.Body, &todo); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "invalid request body")
		return
	}

	if err := utils.ValidateStruct(todo); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "validation failed")
		return
	}

	userContext := middlewares.UserContext(r)
	userID := userContext.UserID

	err = database.UpdateTodoById(*todo.Title, *todo.Body, *todo.Status, *todo.ValidTill, todoID, userID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "invalid payload")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "todo updated successfully",
	})
}

func DeleteTodoByID(w http.ResponseWriter, r *http.Request) {
	userContext := middlewares.UserContext(r)
	userID := userContext.UserID

	todoID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "invalid todo id")
		return
	}

	affected, err := database.DeleteTodoByID(userID, todoID)
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
