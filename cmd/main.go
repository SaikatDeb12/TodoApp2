package main

import (
	"fmt"
	"log"
	"net/http"

)
func main(){
	database.Connect()
	r := routes.SetupRouter()
	fmt.Println("Server running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
