package utils

import (
	"context"
	"errors"

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
