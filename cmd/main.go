package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Saikatdeb12/TodoApp2/internal/database"
	"github.com/Saikatdeb12/TodoApp2/internal/routes"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load .env")
	}

	err := database.Connect()
	if err != nil {
		log.Fatal("failed to connect to Db")
	}
	r := routes.SetupRouter()
	fmt.Println("Server running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
