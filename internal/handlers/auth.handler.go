package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/Saikatdeb12/TodoApp2/internal/database"
	middlewares "github.com/Saikatdeb12/TodoApp2/internal/middleware"
	"github.com/Saikatdeb12/TodoApp2/internal/models"
	"github.com/Saikatdeb12/TodoApp2/internal/utils"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest

	if err := utils.ParseBody(r.Body, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "invalid request payload")
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "validation error")
		return
	}

	isEmailExists, err := database.CheckUserExistsByEmail(req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if isEmailExists {
		utils.RespondJSON(w, http.StatusConflict, map[string]string{
			"error": "user already exists",
		})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "password hashing failed", http.StatusInternalServerError)
		return
	}

	userID, err := database.CreateUser(req.Name, req.Email, hashedPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sessionID, err := database.CreateSession(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.RespondJSON(w, http.StatusCreated, map[string]string{
		"message": "user registered successfully",
		"session": sessionID,
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := utils.ParseBody(r.Body, &req); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, map[string]string{
			"validation_error": err.Error(),
		})
		return
	}

	user, err := database.GetUserAuthByEmail(req.Email)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := utils.CheckPassword(user.Password, req.Password); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	userID := user.UserID
	sessionID, err := database.CreateSession(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.RespondJSON(w, http.StatusCreated, map[string]string{
		"message": "user logged in successfully",
		"session": sessionID,
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	userContext := middlewares.UserContext(r)
	sessionID := userContext.SessionID

	affectedRows, err := database.ArchiveSession(sessionID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "no session found")
		return
	}

	if affectedRows == 0 {
		utils.RespondError(w, http.StatusUnauthorized, err, "invalid session")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})
}
