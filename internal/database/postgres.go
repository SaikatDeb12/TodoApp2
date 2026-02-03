package database

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func GoDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv(key)
}

func Connect() {
	connStr := GoDotEnvVariable("POSTGRESQL_URL")
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
