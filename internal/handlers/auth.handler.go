package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Saikatdeb12/TodoApp2/internal/database"
	"github.com/Saikatdeb12/TodoApp2/internal/models"
	"github.com/Saikatdeb12/TodoApp2/internal/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)


func Register(w http.ResponseWriter, r *http.Request){
	var req models.RegisterRequest
	if parseErr := utils.ParseBody(r.Body, &req); parseErr!=nil {
		
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err!=nil {
		http.Error(w, "Password hashing failed!", http.StatusInternalServerError)
		return
	}

	query := `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
	`

	_,err = database.DB.Exec(query, req.Name, req.Email, string(hash))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string {
		"msg" : "User registered successfully",
	})
}

func Login(w http.ResponseWriter, r *http.Request){
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	var user models.User

	query := `SELECT id, password FROM users WHERE email=$1`
	err := database.DB.QueryRow(query, req.Email).Scan(&user.UserID, &user.Password)
	if err == sql.ErrNoRows{
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	} else if err!=nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err!=nil{
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}

	//create a session:
	sessionID := uuid.New()
	expires := time.Now().Add(24*time.Hour)

	query=`
		INSERT INTO sessions (id, user_id, created_at, expires_at)
		VALUES($1, $2, now(), $3)
	`

	_, err = database.DB.Exec(query, sessionID, user.UserID, expires)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string] string {
		"session" : sessionID.String(),
	})
}

func Logout(w http.ResponseWriter, r *http.Request){
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

	query := `DELETE FROM sessions WHERE id=$1`
	res, err := database.DB.Exec(query, sessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := res.RowsAffected()
	if affected==0 {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string] string {
		"msg" : "Logged out successfully",
	})

}
