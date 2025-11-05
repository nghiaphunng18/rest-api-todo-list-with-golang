package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Starting server on port 8080")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Todo App")
	})

	http.ListenAndServe(":8080", nil)
}