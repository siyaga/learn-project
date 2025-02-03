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


    // Product routes
    api.Post("/products", controllers.CreateProduct)   // Create a product
    api.Get("/products", controllers.GetProducts)      // Get all products
    api.Get("/products/:id", controllers.GetProduct)   // Get a single product
    api.Put("/products/:id", controllers.UpdateProduct) // Update a product
    api.Delete("/products/:id", controllers.DeleteProduct) // Delete a product
   
}