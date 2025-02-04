package controllers

import (
	"learn_project/database"
	"learn_project/models"
	"learn_project/utils"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Struct for input validation
type BankInput struct {
	BankName  string `json:"bank_name" validate:"required"`
	AccountNo string `json:"account_no" validate:"required"`
}

// Add a new bank account (CREATE)
func AddBank(c *fiber.Ctx) error {
	email, ok := c.Locals("email").(string)
	if !ok {
		return utils.ResponseError(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	// Get user
	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return utils.ResponseError(c, fiber.StatusNotFound, "User not found", nil)
	}

	var input BankInput
	if err := c.BodyParser(&input); err != nil {
		return utils.ResponseError(c, fiber.StatusBadRequest, "Invalid input", nil)
	}

	// Check if account number is unique
	var existingBank models.Bank
	if err := database.DB.Where("account_no = ?", input.AccountNo).First(&existingBank).Error; err == nil {
		return utils.ResponseError(c, fiber.StatusConflict, "Account number already in use", nil)
	}

	bank := models.Bank{
		UserID:    user.ID,
		BankName:  input.BankName,
		AccountNo: input.AccountNo,
		Nominal:   0, // Default balance 0
	}

	if err := database.DB.Create(&bank).Error; err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Could not add bank", nil)
	}

	return utils.ResponseSuccessOneData(c, "Bank added successfully", fiber.Map{
		"id":         bank.ID,
		"bank_name":  bank.BankName,
		"account_no": bank.AccountNo,
		"nominal":    bank.Nominal,
	})
}

// Get all banks for a user (READ)
func GetUserBanks(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	search := c.Query("search", "")
	offset := (page - 1) * limit

	email, ok := c.Locals("email").(string)
	if !ok {
		return utils.ResponseError(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return utils.ResponseError(c, fiber.StatusNotFound, "User not found", nil)
	}

	var banks []models.Bank
	query := database.DB.Where("user_id = ?", user.ID).Offset(offset).Limit(limit)

	if search != "" {
		search = strings.ToLower(search)
		query = query.Where("LOWER(bank_name) LIKE ? OR account_no LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Find(&banks).Error; err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Could not fetch banks", nil)
	}

	var count int64
	countQuery := database.DB.Model(&models.Bank{}).Where("user_id = ?", user.ID)

	if search != "" {
		countQuery = countQuery.Where("LOWER(bank_name) LIKE ? OR account_no LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := countQuery.Count(&count).Error; err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Could not fetch bank count", nil)
	}

	return utils.ResponseSuccessManyData(c, "Banks retrieved successfully", banks, page, limit, int(count))
}

// Update bank details (UPDATE)
func UpdateBank(c *fiber.Ctx) error {
	bankID := c.Params("id")

	var bank models.Bank
	if err := database.DB.First(&bank, "id = ?", bankID).Error; err != nil {
		return utils.ResponseError(c, fiber.StatusNotFound, "Bank not found", nil)
	}

	var input BankInput
	if err := c.BodyParser(&input); err != nil {
		return utils.ResponseError(c, fiber.StatusBadRequest, "Invalid input", nil)
	}

	bank.BankName = input.BankName
	bank.AccountNo = input.AccountNo

	if err := database.DB.Save(&bank).Error; err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Could not update bank", nil)
	}

	return utils.ResponseSuccessOneData(c, "Bank updated successfully", bank)
}

// Delete bank (DELETE)
func DeleteBank(c *fiber.Ctx) error {
	bankID := c.Params("id")

	if err := database.DB.Delete(&models.Bank{}, "id = ?", bankID).Error; err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Could not delete bank", nil)
	}

	return utils.ResponseSuccessOneData(c, "Bank deleted successfully",nil)
}

// Add money to bank (UPDATE Nominal)
type AddMoneyInput struct {
	Amount float64 `json:"amount" validate:"required,min=1"`
}

func AddMoney(c *fiber.Ctx) error {
	bankID := c.Params("id")

	var bank models.Bank
	if err := database.DB.First(&bank, "id = ?", bankID).Error; err != nil {
		return utils.ResponseError(c, fiber.StatusNotFound, "Bank not found", nil)
	}

	var input AddMoneyInput
	if err := c.BodyParser(&input); err != nil {
		return utils.ResponseError(c, fiber.StatusBadRequest, "Invalid input", nil)
	}

	// Update nominal balance
	bank.Nominal += input.Amount

	if err := database.DB.Save(&bank).Error; err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Could not update balance", nil)
	}

	return utils.ResponseSuccessOneData(c, "Money added successfully", fiber.Map{
		"id":         bank.ID,
		"bank_name":  bank.BankName,
		"account_no": bank.AccountNo,
		"nominal":    bank.Nominal,
	})
}
