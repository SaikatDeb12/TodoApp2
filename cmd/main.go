package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Saikatdeb12/TodoApp2/internal/database"
	"github.com/Saikatdeb12/TodoApp2/internal/routes"
	"github.com/Saikatdeb12/TodoApp2/internal/utils"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load .env")
	}

	err := database.Connect()
	if err != nil {
		log.Fatal("failed to connect to DB")
	}

	r := routes.SetupRouter()
	serverPort := utils.GoDotEnvVariable("SERVER_PORT")
	addr := fmt.Sprintf(":%s", serverPort)
	fmt.Printf("Server running on port %s", serverPort)
	log.Fatal(http.ListenAndServe(addr, r))
}
