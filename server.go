package main

import (
	"learn_project/database"
	"learn_project/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
    // Load environment variables from .env file
    if err := godotenv.Load(); err != nil {
        panic("Error loading .env file")
    }

    // Initialize the database
    database.Connect()

    // Create a new Fiber app
    app := fiber.New()

    // Set up routes
    routes.SetupRoutes(app)

    // Start the server
    app.Listen(":3000")
}