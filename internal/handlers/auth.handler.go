package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	dbhelper "github.com/Saikatdeb12/TodoApp2/internal/database/dbHelper"
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

	isEmailExists, err := dbhelper.CheckUserExistsByEmail(req.Email)
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

	userID, err := dbhelper.CreateUser(req.Name, req.Email, hashedPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sessionID, err := dbhelper.CreateSession(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := utils.GenerateJWT(userID, sessionID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "error while generating token")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, map[string]string{
		"message": "user registered successfully",
		"token":   token,
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

	user, err := dbhelper.GetUserAuthByEmail(req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.RespondError(w, http.StatusUnauthorized, err, "invalid credentials")
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, err, "no users fetched")
	}

	if err := utils.CheckPassword(user.Password, req.Password); err != nil {
		utils.RespondError(w, http.StatusUnauthorized, err, "invalid credentials")
		return
	}

	userID := user.UserID

	sessionID, err := dbhelper.CreateSession(userID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to create session")
		return
	}

	token, err := utils.GenerateJWT(userID, sessionID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "error while generating token")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, map[string]string{
		"message": "user logged in successfully",
		"token":   token,
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	userContext := middlewares.UserContext(r)
	sessionID := userContext.SessionID

	affectedRows, err := dbhelper.ArchiveSession(sessionID)
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
