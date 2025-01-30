// controllers/auth.go
package controllers

import (
	"learn_project/database"
	"learn_project/models"
	"learn_project/utils"

	"github.com/gofiber/fiber/v2"
)

// Struct untuk request body Register
type RegisterInput struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func Register(c *fiber.Ctx) error {
	var input RegisterInput
	if err := c.BodyParser(&input); err != nil {
		return utils.ResponseError(c, fiber.StatusBadRequest, "Invalid input", nil)
	}

	if input.Name == "" || input.Email == "" || input.Password == "" {
		return utils.ResponseError(c, fiber.StatusBadRequest, "All fields are required", nil)
	}

	var existingUser models.User
	if err := database.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		return utils.ResponseError(c, fiber.StatusConflict, "Email already in use", nil)
	}

	user := models.User{
		Name:  input.Name,
		Email: input.Email,
	}

	if err := user.HashPassword(input.Password); err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Could not hash password", nil)
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Could not create user", nil)
	}

	return utils.ResponseSuccessOneData(c, "User registered successfully", fiber.Map{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}

func Login(c *fiber.Ctx) error {
	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return utils.ResponseError(c, fiber.StatusBadRequest, "Invalid input", nil)
	}

	// Cari user berdasarkan email
	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return utils.ResponseError(c, fiber.StatusUnauthorized, "Invalid credentials", nil)
	}

	// Periksa password
	if err := user.CheckPassword(input.Password); err != nil {
		return utils.ResponseError(c, fiber.StatusUnauthorized, "Invalid credentials", nil)
	}

	// Generate JWT Token
	accessToken, exp, err := utils.GenerateToken(user.Email)
	if err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Could not generate token", nil)
	}

	// Generate Refresh Token
	refreshToken, err := utils.GenerateRefreshToken(user.Email)
	if err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Could not generate refresh token", nil)
	}

	// Return response dengan user data dan token
	return utils.ResponseSuccessOneData(c, "Login successful", fiber.Map{
		"user": fiber.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
		"access_token":  accessToken,
		"expires_at":    exp,  // Waktu kadaluarsa token
		"refresh_token": refreshToken,
	})
}

func GetUser(c *fiber.Ctx) error {
	email, ok := c.Locals("email").(string)
	if !ok {
		return utils.ResponseError(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return utils.ResponseError(c, fiber.StatusNotFound, "User not found", nil)
	}

	return utils.ResponseSuccessOneData(c, "User data retrieved successfully", fiber.Map{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}
