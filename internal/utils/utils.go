package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Saikatdeb12/TodoApp2/internal/models"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var (
	validate  = validator.New()
	SecretKey = GoDotEnvVariable("JWT_SECRET_KEY")
)

func GoDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("error loading .env file")
		return ""
	}

	return os.Getenv(key)
}

func ValidateStruct(payload interface{}) error {
	return validate.Struct(payload)
}

func ParseBody(body io.Reader, out interface{}) error {
	return json.NewDecoder(body).Decode(out)
}

func EncodeBody(res http.ResponseWriter, data interface{}) error {
	return json.NewEncoder(res).Encode(data)
}

func RespondJSON(w http.ResponseWriter, statusCode int, body interface{}) {
	w.WriteHeader(statusCode)
	if body != nil {
		if err := EncodeBody(w, body); err != nil {
			log.Printf("failed with error: %+v", err)
		}
	}
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func RespondError(w http.ResponseWriter, statusCode int, err error, message string) {
	w.WriteHeader(statusCode)
	var errStr string

	if err != nil {
		errStr = err.Error()
	}

	NewError := models.Error{
		Error:      errStr,
		StatusCode: statusCode,
		Message:    message,
	}

	if err := EncodeBody(w, NewError); err != nil {
		fmt.Printf("error is %v", err)
	}
}

func ParseDate(dateStr string) (*time.Time, error) {
	if dateStr == "" {
		return nil, nil
	}
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, err
	}
	return &parsedDate, nil
}

func GenerateJWT(userID, sessionID string) (string, error) {
	claims := jwt.MapClaims{
		"userID":    userID,
		"sessionID": sessionID,
		"expiresAt": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString([]byte(SecretKey))
}
