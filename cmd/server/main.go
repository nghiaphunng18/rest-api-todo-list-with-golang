package main

import (
	"fmt"
	"log"
	"net/http"
	"todo-api/internal/config"
	"todo-api/internal/database"
)

func main() {
	// load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("cannot load configuration: %v", err)
	}
	
	// connect to database
	_, err = database.InitDB(cfg)
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}
	log.Println("Connected to database success")

	// run server
	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	fmt.Printf("Starting server on port %s\n", cfg.ServerPort)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Todo App")
	})

	log.Fatal(http.ListenAndServe(serverAddr, nil))
}
