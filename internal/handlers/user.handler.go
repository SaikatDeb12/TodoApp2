package handlers

import (
	"net/http"

	"github.com/Saikatdeb12/TodoApp2/internal/database"
	dbhelper "github.com/Saikatdeb12/TodoApp2/internal/database/dbHelper"
	"github.com/Saikatdeb12/TodoApp2/internal/middleware"
	"github.com/Saikatdeb12/TodoApp2/internal/utils"
	"github.com/jmoiron/sqlx"
)

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userCtx, ok := middleware.UserContext(r)

	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
		return
	}

	userID := userCtx.UserID
	sessionID := userCtx.SessionID

	trxErr := database.Tx(func(tx *sqlx.Tx) error {
		if TodoErr := dbhelper.DeleteAllTodosTx(tx, userID); TodoErr != nil {
			return TodoErr
		}

		if TodoErr := dbhelper.ArchiveSessionTx(tx, sessionID); TodoErr != nil {
			return TodoErr
		}

		return dbhelper.ArchiveUserTx(tx, userID)
	})

	if trxErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, trxErr, "user deletion failed")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "user deleted successfully",
	})
}
