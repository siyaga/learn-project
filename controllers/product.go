// controllers/product.go
package controllers

import (
	"learn_project/database"
	"learn_project/models"
	"learn_project/utils"

	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Struct untuk request body CreateProduct
type CreateProductInput struct {
    Name        string  `json:"name" validate:"required"`
    Description string  `json:"description"`
    Price       float64 `json:"price" validate:"required,min=0"`
}

// Struct untuk request body UpdateProduct
type UpdateProductInput struct {
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price" validate:"min=0"`
}

// CreateProduct creates a new product
func CreateProduct(c *fiber.Ctx) error {
    var input CreateProductInput
    if err := c.BodyParser(&input); err != nil {
        return utils.ResponseError(c, fiber.StatusBadRequest, "Invalid input", nil)
    }

    if input.Name == "" || input.Price <= 0 {
        return utils.ResponseError(c, fiber.StatusBadRequest, "Name and price are required", nil)
    }

    product := models.Product{
        Name:        input.Name,
        Description: input.Description,
        Price:       input.Price,
    }

    if err := database.DB.Create(&product).Error; err != nil {
        return utils.ResponseError(c, fiber.StatusInternalServerError, "Could not create product", nil)
    }

    return utils.ResponseSuccessOneData(c, "Product created successfully", fiber.Map{
        "id":          product.ID,
        "name":        product.Name,
        "description": product.Description,
        "price":       product.Price,
    })
}

func GetProducts(c *fiber.Ctx) error {
	// Parse query parameters for pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))    // Default page is 1
	limit, _ := strconv.Atoi(c.Query("limit", "10")) // Default limit is 10

	// Get the search term from query parameters
	search := c.Query("search", "") // Default is an empty string

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch products from the database with pagination and search
	var products []models.Product
	query := database.DB.Offset(offset).Limit(limit)

	// Add search filter if search term is provided
	if search != "" {
			search = strings.ToLower(search) // Convert search term to lowercase for case-insensitive search
			query = query.Where("LOWER(name) LIKE ?", "%"+search+"%")
	}

	if err := query.Find(&products).Error; err != nil {
			return utils.ResponseError(c, fiber.StatusInternalServerError, "Could not fetch products", nil)
	}

	// Get the total count of products (with search filter if applicable)
	var count int64
	countQuery := database.DB.Model(&models.Product{})

	if search != "" {
			countQuery = countQuery.Where("LOWER(name) LIKE ?", "%"+search+"%")
	}

	if err := countQuery.Count(&count).Error; err != nil {
			return utils.ResponseError(c, fiber.StatusInternalServerError, "Could not fetch product count", nil)
	}

	// Return the response with pagination details
	return utils.ResponseSuccessManyData(c, "Products retrieved successfully", products, page, limit, int(count))
}


// GetProduct fetches a single product by ID
func GetProduct(c *fiber.Ctx) error {
    id := c.Params("id")
    var product models.Product
    if err := database.DB.Where("id = ?", id).First(&product).Error; err != nil {
        return utils.ResponseError(c, fiber.StatusNotFound, "Product not found", nil)
    }

    return utils.ResponseSuccessOneData(c, "Product retrieved successfully", fiber.Map{
        "id":          product.ID,
        "name":        product.Name,
        "description": product.Description,
        "price":       product.Price,
    })
}

// UpdateProduct updates a product by ID
func UpdateProduct(c *fiber.Ctx) error {
    id := c.Params("id")
    var product models.Product
    if err := database.DB.Where("id = ?", id).First(&product).Error; err != nil {
        return utils.ResponseError(c, fiber.StatusNotFound, "Product not found", nil)
    }

    var input UpdateProductInput
    if err := c.BodyParser(&input); err != nil {
        return utils.ResponseError(c, fiber.StatusBadRequest, "Invalid input", nil)
    }

    if input.Name != "" {
        product.Name = input.Name
    }
    if input.Description != "" {
        product.Description = input.Description
    }
    if input.Price > 0 {
        product.Price = input.Price
    }

    if err := database.DB.Save(&product).Error; err != nil {
        return utils.ResponseError(c, fiber.StatusInternalServerError, "Could not update product", nil)
    }

    return utils.ResponseSuccessOneData(c, "Product updated successfully", fiber.Map{
        "id":          product.ID,
        "name":        product.Name,
        "description": product.Description,
        "price":       product.Price,
    })
}

// DeleteProduct deletes a product by ID
func DeleteProduct(c *fiber.Ctx) error {
    id := c.Params("id")
    var product models.Product
    if err := database.DB.Where("id = ?", id).First(&product).Error; err != nil {
        return utils.ResponseError(c, fiber.StatusNotFound, "Product not found", nil)
    }

    if err := database.DB.Delete(&product).Error; err != nil {
        return utils.ResponseError(c, fiber.StatusInternalServerError, "Could not delete product", nil)
    }

    return utils.ResponseSuccessOneData(c, "Product deleted successfully", nil)
}