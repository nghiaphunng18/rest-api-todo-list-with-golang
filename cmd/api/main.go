package main

import (
	"fmt"
	"log"
	"net/http"
	"todo-api/internal/database"
)

func main() {
	// connect to database
	_, err := database.InitDB()
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}
	log.Println("Connected to database success")

	// run server
	fmt.Println("Starting server on port 8080")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Todo App")
	})

	http.ListenAndServe(":8080", nil)
}