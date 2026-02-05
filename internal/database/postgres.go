package database

import (
	"fmt"

	"github.com/Saikatdeb12/TodoApp2/internal/utils"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func Connect() {
	connStr := utils.GoDotEnvVariable("POSTGRESQL_URL")
	var err error
	DB, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		panic("DB open failed!")
	}

	if err = DB.Ping(); err != nil {
		panic("DB ping error")
	}

	fmt.Println("Connected to PostgreSQL")
}
