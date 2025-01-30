package routes

import (
	"learn_project/controllers"
	"learn_project/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
    // Public routes (no authentication required)
    app.Post("/register", controllers.Register) // Register a new user
    app.Post("/login", controllers.Login)       // Login and get JWT token

    // Protected routes (require JWT authentication)
    api := app.Group("/api", middleware.Protected()) // Group for protected routes
    api.Get("/user", controllers.GetUser)            // Example protected route
}