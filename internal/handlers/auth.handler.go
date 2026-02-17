package handlers

import (
	"database/sql"
	"errors"
	"fmt"
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

	err := dbhelper.CheckUserExistsByEmail(req.Email)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "email already exists")
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "password hasing failed")
		return
	}

	userID, err := dbhelper.CreateUser(req.Name, req.Email, hashedPassword)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to create user")
		return
	}

	sessionID, err := dbhelper.CreateSession(userID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to create session")
		return
	}

	token, err := utils.GenerateJWT(userID, sessionID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "token not created")
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
		utils.RespondError(w, http.StatusBadRequest, err, "invalid payload")
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
		return
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
	userContext, ok := middlewares.UserContext(r)
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
	}
	sessionID := userContext.SessionID

	fmt.Print(sessionID)
	err := dbhelper.ValidateUserSession(sessionID)
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, err, "no active session")
		return
	}

	err = dbhelper.ArchiveSession(sessionID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "no session found")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})
}
