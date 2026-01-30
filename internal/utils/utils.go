package utils

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	middlewares "github.com/Saikatdeb12/TodoApp2/internal/middleware"
	"github.com/google/uuid"
)


func GetUserID(ctx context.Context) (uuid.UUID, error) {
	userId, ok := ctx.Value(middlewares.UserIDkey).(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("Unauthorized")
	}

	return userId,nil
}

func ParseBody(body io.Reader, out interface{}) error{
	return json.NewDecoder(body).Decode(out)
}

func EncodeBody(res http.ResponseWriter, data interface{}) error {
	return json.NewEncoder(res).Encode(data)
}

func RespondJSON(w http.ResponseWriter, statusCode int, body interface{}){
	w.WriteHeader(statusCode)
	if body!=nil {
		if err:= EncodeBody(w, body); err != nil {
			log.Fatal("Failed to respond JSON with error: %+v", err)
		}
	}
}



// func ResponseError(w http.ResponseWriter, statusCode int, err error, messageToUser string ){
// 	log.Fatal()
// }
