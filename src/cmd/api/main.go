package main

import (
	"log"
	"net/http"

	"users-api/src/internal/config"
	"users-api/src/internal/db"
	"users-api/src/internal/delivery/handlers"
	httpDelivery "users-api/src/internal/delivery/http"
	"users-api/src/internal/repository/postgres"
	"users-api/src/internal/service"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	database, err := db.NewDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	userRepo := postgres.NewUserRepository(database.DB)

	userService := service.NewUserService(userRepo)

	userHandler := handlers.NewUserHandler(userService)

	router := httpDelivery.NewRouter(userHandler)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
