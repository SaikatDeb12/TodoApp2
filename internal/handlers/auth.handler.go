package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/Saikatdeb12/TodoApp2/internal/database"
	middlewares "github.com/Saikatdeb12/TodoApp2/internal/middleware"
	"github.com/Saikatdeb12/TodoApp2/internal/models"
	"github.com/Saikatdeb12/TodoApp2/internal/utils"
	"github.com/google/uuid"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest

	if err := utils.ParseBody(r.Body, &req); err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request payload",
		})
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		// send actual error
		utils.RespondJSON(w, http.StatusBadRequest, map[string]string{
			"validation_error": err.Error(),
		})
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

	sessionID := uuid.New()
	sessionExpiresAt := time.Now().Add(24 * time.Hour)

	if err := database.CreateSession(sessionID, userID, sessionExpiresAt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.RespondJSON(w, http.StatusCreated, map[string]string{
		"message": "user registered successfully",
		"session": sessionID.String(),
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

	sessionID := uuid.New()
	sessionExpiresAt := time.Now().Add(24 * time.Hour)

	if err := database.CreateSession(sessionID, user.UserID, sessionExpiresAt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"session": sessionID.String(),
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	userContext := middlewares.UserContext(r)
	sessionID := userContext.SessionID

	// no need of affected rows
	affectedRows, err := database.ArchiveSession(sessionID)
	if err != nil {
		// utils respond error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if affectedRows == 0 {
		http.Error(w, "invalid session", http.StatusUnauthorized)
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})
}
