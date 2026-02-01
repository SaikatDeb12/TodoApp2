package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/Saikatdeb12/TodoApp2/internal/database"
	"github.com/Saikatdeb12/TodoApp2/internal/models"
	"github.com/Saikatdeb12/TodoApp2/internal/utils"
	"github.com/google/uuid"
)

// {"session":"066dfad3-fadd-451c-b91a-a4b5fe73f30e"}

func Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest

	if err := utils.ParseBody(r.Body, &req); err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload",
		})
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, map[string]string{
			"validation_error": err.Error(),
		})
		return
	}

	exists, err := database.CheckUserExistsByEmail(req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if exists {
		utils.RespondJSON(w, http.StatusConflict, map[string]string{
			"error": "User already exists",
		})
		return
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Password hashing failed", http.StatusInternalServerError)
		return
	}

	if err := database.CreateUser(req.Name, req.Email, hash); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.RespondJSON(w, http.StatusCreated, map[string]string{
		"msg": "User registered successfully",
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest

	if err := utils.ParseBody(r.Body, &req); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, map[string]string{
			"validation_error": err.Error(),
		})
		return
	}

	user, err := database.GetUserAuthByEmail(req.Email)
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := utils.CheckPassword(user.Password, req.Password); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	sessionID := uuid.New()
	expires := time.Now().Add(24 * time.Hour)

	if err := database.CreateSession(sessionID, user.UserID, expires); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"session": sessionID.String(),
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	sessionID, err := uuid.Parse(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	affected, err := database.DeleteSession(sessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if affected == 0 {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"msg": "Logged out successfully",
	})
}
