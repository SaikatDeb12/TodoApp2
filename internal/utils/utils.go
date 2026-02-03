package utils

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	middlewares "github.com/Saikatdeb12/TodoApp2/internal/middleware"
	"github.com/Saikatdeb12/TodoApp2/internal/models"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

func ValidateStruct(payload interface{}) error {
	return validate.Struct(payload)
}

func GetUserID(ctx context.Context) (uuid.UUID, error) {
	requestContext, err := ctx.Value(middlewares.RequestContextKey).(models.RequestContext)
	if !err {
		return uuid.Nil, errors.New("Unauthorized")
	}

	return requestContext.UserID, nil
}

func GetSessionID(ctx context.Context) (uuid.UUID, error) {
	requestContext, err := ctx.Value(middlewares.RequestContextKey).(models.RequestContext)
	if !err {
		return uuid.Nil, errors.New("Unauthorized")
	}

	return requestContext.SessionID, nil
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
			log.Printf("Failed with error: %+v", err)
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
