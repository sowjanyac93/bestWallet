package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sowjanyac93/bestWallet/modules/initializers"
	"github.com/sowjanyac93/bestWallet/modules/models"
)

// GetAccounts fetches all accounts
func GetAccounts(c *gin.Context) {
	accounts := []models.Account{}
	result := initializers.DB.Find(&accounts)
	if result.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"Error": "Failed to fetch accounts",
		})
	} else {
		// Marshal the accounts using custom MarshalJSON
		serializedAccounts, err := json.Marshal(accounts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Error": "Failed to marshal accounts",
			})
			return
		}

		// Respond with the marshaled data
		c.Data(http.StatusOK, "application/json", serializedAccounts)
	}
}

// GetPersonalAccount fetches the accounts for the authenticated user
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

	// Marshal the accounts using custom MarshalJSON
	serializedAccounts, err := json.Marshal(accounts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Failed to marshal accounts",
		})
		return
	}

	// Respond with the marshaled data
	c.Data(http.StatusOK, "application/json", serializedAccounts)
}

// AddAccountForUserId adds an account for the specified user ID
func AddAccountForUserId(c *gin.Context) {
	// Retrieve the userId from the URL parameters
	userID, err := uuid.Parse(c.Param("id"))
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
	account.UserId = userID

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

	// Marshal the account using custom MarshalJSON
	serializedAccount, err := json.Marshal(account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Failed to marshal account",
		})
		return
	}

	// Respond with the marshaled data
	c.Data(http.StatusOK, "application/json", serializedAccount)
}

// AddSelfAccount adds an account for the authenticated user
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
	account.UserId = userID

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

	// Marshal the account using custom MarshalJSON
	serializedAccount, err := json.Marshal(account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Failed to marshal account",
		})
		return
	}

	// Respond with the marshaled data
	c.Data(http.StatusOK, "application/json", serializedAccount)
}

// DeleteAccount deletes an account
func DeleteAccount(c *gin.Context) {
	accountID := c.Param("id")

	// Attempt to delete the account
	result := initializers.DB.Where("id = ?", accountID).Delete(&models.Account{})

	if result.Error != nil {
		// Database error
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"Error": "Failed to delete account",
		})
		return
	}

	if result.RowsAffected == 0 {
		// Account not found
		c.JSON(http.StatusNotFound, gin.H{
			"Error": "Account not found",
		})
		return
	}

	// Account deleted successfully
	c.JSON(http.StatusOK, gin.H{
		"message": "Account deleted successfully",
	})
}

// Fetches details for particular AccountID
func GetAccountDetails(c *gin.Context) {
	accountID := c.Param("id")

	// Fetch the account with the specified accountId
	var account models.Account
	result := initializers.DB.Where("id = ?", accountID).First(&account)

	if result.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"Error": "Failed to fetch account details",
		})
		return
	}

	// Marshal the account using custom MarshalJSON
	serializedAccount, err := account.MarshalJSON()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Failed to marshal account details",
		})
		return
	}

	// Respond with the marshaled data
	c.Data(http.StatusOK, "application/json", serializedAccount)
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
