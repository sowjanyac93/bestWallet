package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sowjanyac93/bestWallet/modules/initializers"
	"github.com/sowjanyac93/bestWallet/modules/models"
)

func GetAccounts(c *gin.Context) {
	accounts := []models.Account{}
	result := initializers.DB.Find(&accounts)
	if result.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"Error": "Failed to fetch accounts",
		})
	} else {
		c.JSON(http.StatusOK, &accounts)
	}
}

func GetPersonalAccount(c *gin.Context) {
	// Retrieve the user from the Gin context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User information not available"})
		return
	}

	// Extract the user ID
	userID := user.(models.User).ID

	var accounts []models.Account
	result := initializers.DB.Find(&accounts, "user_id = ?", userID)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch account data"})
		return
	}

	var accountData []gin.H
	for _, acc := range accounts {
		accountData = append(accountData, gin.H{
			"accountNo": acc.AccountNumber,
			"bankName":  acc.BankName,
			"balance":   acc.Balance,
			"status":    acc.Status,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"account_data": accountData,
	})
}

func AddAccountForUserId(c *gin.Context) {
	// Retrieve the userId from the URL parameters
	userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Invalid user ID",
		})
		return
	}

	var account models.Account
	if err := c.BindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read input data",
		})
		return
	}

	// Set the user ID for the account
	account.UserId = int(userId)

	// Create the account in the database
	result := initializers.DB.Create(&account)
	if result.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"Error": "Failed to create account",
		})
		return
	}

	// Simulate KYC process
	go simulateKYCProcess(account.ID, c)

	c.JSON(http.StatusOK, gin.H{
		"accountNo": account.AccountNumber,
		"bankName":  account.BankName,
		"balance":   account.Balance,
		"message":   "Account created successfully",
	})
}

// Function to add an account for the authenticated user
func AddSelfAccount(c *gin.Context) {
	// Retrieve the user from the Gin context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User information not available"})
		return
	}

	// Extract the user ID
	userID := user.(models.User).ID

	var account models.Account
	if err := c.BindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read input data",
		})
		return
	}

	// Set the user ID for the account
	account.UserId = int(userID)

	// Create the account in the database
	result := initializers.DB.Create(&account)
	if result.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"Error": "Failed to create account",
		})
		return
	}

	// Simulate KYC process
	go simulateKYCProcess(account.ID, c)

	c.JSON(http.StatusOK, gin.H{
		"accountNo": account.AccountNumber,
		"bankName":  account.BankName,
		"balance":   account.Balance,
		"message":   "Account created successfully",
	})
}

func DeleteAccount(c *gin.Context) {
	var account models.Account
	result := initializers.DB.Where("id = ?", c.Param("id")).Delete(&account)
	if result.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"Error": "Failed to delete account",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "account deleted successfully",
		})
	}
}

// Simulate KYC process
func simulateKYCProcess(accountId uint, c *gin.Context) {
	// Sleep for 60 seconds
	time.Sleep(60 * time.Second)

	tx := initializers.DB.Begin()
	if err := tx.Model(&models.Account{}).Where("id = ?", accountId).Update("status", "Approved").Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"Error": "KYC failed",
		})
		return
	}
	tx.Commit()
}
