package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Saikatdeb12/TodoApp2/internal/database"
	"github.com/Saikatdeb12/TodoApp2/internal/routes"
)
func main(){
	database.Connect()
	r := routes.SetupRouter()
	fmt.Println("Server running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
