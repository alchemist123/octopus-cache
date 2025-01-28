package main

import (
	"fmt"
	"net/http"

	"octopus-cache/internal/database"
	"octopus-cache/pkg/api"
)

func main() {
	db := database.NewDatabase()
	handler := api.NewHandler(db)

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", handler)
}
