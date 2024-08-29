package main

import (
	"Blogsite/config"
	"Blogsite/routes"
	"log"
	"net/http"
)

func main() {
	config.InitDB()

	// Set up the router
	router := routes.SetupRoutes()

	// Start the server
	log.Println("Server running on port 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
