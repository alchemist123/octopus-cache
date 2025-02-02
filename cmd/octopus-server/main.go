package main

import (
	"fmt"
	"net/http"

	"octopus-cache/internal/database"
	"octopus-cache/pkg/api"
)

func main() {
	db, err := database.NewDatabase("defaultDataDir")
	if err != nil {
		fmt.Printf("failed to initialize database: %v\n", err)
		return
	}
	handler := api.NewHandler(db)

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", handler)
}
