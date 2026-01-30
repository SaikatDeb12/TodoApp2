package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB
func GoDotEnvVariable(key string) string{
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv(key)
}

func Connect(){
	connStr := GoDotEnvVariable("POSTGRESQL_URL")
	var err error
	DB, err = sql.Open("postgres", connStr)
	if(err != nil){
		panic("DB open failed!")
	}

	if err=DB.Ping(); err!=nil{
		panic("DB ping error")
	}

	fmt.Println("Connected to PostgreSQL")
}
