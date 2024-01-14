package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sowjanyac93/bestWallet/modules/initializers"
	"github.com/sowjanyac93/bestWallet/modules/models"
)

func DepositAmount(c *gin.Context) {
	var body struct {
		AccountNumber int
		Amount        int
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read input data",
		})
		return
	}

	var account models.Account
	accDetails := initializers.DB.First(&account, "account_number = ?", body.AccountNumber)

	if accDetails.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Account number does not exist",
		})
		return
	}

	if account.Status == "Pending" {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "KYC not verified",
		})
		return
	}

	// Fetch the userId from the account
	userId := account.UserId

	// Start a database transaction
	tx := initializers.DB.Begin()

	// Update the account balance
	account.Balance += body.Amount
	if err := tx.Save(&account).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Failed to update account balance",
		})
		return
	}

	// Insert a record in the Transaction table
	transaction := models.Transaction{
		UserId:    userId,
		AccountId: int(account.ID),
		Flow:      "Deposit",
		Amount:    body.Amount,
		Details:   "Deposit operation",
	}
	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Failed to create a transaction record",
		})
		return
	}

	// Commit the transaction
	tx.Commit()

	// Simulate KYT process
	go simulateKYTProcess(transaction.ID, c)

	c.JSON(http.StatusOK, gin.H{
		"Message": "Deposit operation completed successfully",
	})
}

func GetTransactions(c *gin.Context) {
	transactions := []models.Transaction{}
	result := initializers.DB.Find(&transactions)
	if result.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"Error": "Failed to fetch transactions",
		})
	} else {
		c.JSON(http.StatusOK, &transactions)
	}
}

func WithdrawAmount(c *gin.Context) {
	var body struct {
		FromAccountNumber int
		Amount            int
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read input data",
		})
		return
	}

	var account models.Account
	accDetails := initializers.DB.First(&account, "account_number = ?", body.FromAccountNumber)

	if accDetails.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Account number does not exist",
		})
		return
	}

	if account.Status == "Pending" {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "KYC not verified",
		})
		return
	}

	// Fetch the userId from the account
	userId := account.UserId

	// Start a database transaction
	tx := initializers.DB.Begin()

	if account.Balance < body.Amount {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Insufficient balance",
		})
		return
	}

	// Update the account balance
	account.Balance -= body.Amount
	if err := tx.Save(&account).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Failed to update account balance",
		})
		return
	}

	// Insert a record in the Transaction table
	transaction := models.Transaction{
		UserId:    userId,
		AccountId: int(account.ID),
		Flow:      "Withdraw",
		Amount:    body.Amount,
		Details:   "Withdrawal operation",
	}
	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Failed to create a transaction record",
		})
		return
	}

	// Commit the transaction
	tx.Commit()

	// Simulate KYT process
	go simulateKYTProcess(transaction.ID, c)

	c.JSON(http.StatusOK, gin.H{
		"Message": "Withdraw operation completed successfully",
	})
}

func SelfTransfer(c *gin.Context) {
	var body struct {
		FromAccountNumber int
		ToAccountNumber   int
		Amount            int
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read input data",
		})
		return
	}

	var fromAccount models.Account
	fromAccDetails := initializers.DB.First(&fromAccount, "account_number = ?", body.FromAccountNumber)

	if fromAccDetails.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "From Account number does not exist",
		})
		return
	}

	if fromAccount.Status == "Pending" {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "From account KYC not verified",
		})
		return
	}

	var toAccount models.Account
	toAccDetails := initializers.DB.First(&toAccount, "account_number = ?", body.ToAccountNumber)

	if toAccDetails.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "To Account number does not exist",
		})
		return
	}

	if toAccount.Status == "Pending" {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "To account KYC not verified",
		})
		return
	}

	// Fetch the userId from the account
	fromUserId := fromAccount.UserId
	toUserId := toAccount.UserId

	if fromUserId != toUserId {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Self-transfer not possible",
		})
		return
	}

	// Start a database transaction
	tx := initializers.DB.Begin()

	if fromAccount.Balance < body.Amount {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Insufficient balance",
		})
		return
	}

	// Update the FROM account balance
	fromAccount.Balance -= body.Amount
	if err := tx.Save(&fromAccount).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Failed to update FROM Account balance",
		})
		return
	}

	// Update the TO account balance
	toAccount.Balance += body.Amount
	if err := tx.Save(&toAccount).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Failed to update To Account balance",
		})
		return
	}

	// Insert a record in the Transaction table
	transaction := models.Transaction{
		UserId:    fromUserId,
		AccountId: int(fromAccount.ID),
		Flow:      "Self-transfer",
		Amount:    body.Amount,
		Details:   strconv.Itoa(toAccount.AccountNumber),
	}
	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Failed to create a transaction record",
		})
		return
	}

	// Commit the transaction
	tx.Commit()

	// Simulate KYT process
	go simulateKYTProcess(transaction.ID, c)

	c.JSON(http.StatusOK, gin.H{
		"Message": "Self-transfer operation completed successfully",
	})
}

// Simulate KYT process
func simulateKYTProcess(transactionID uint, c *gin.Context) {
	// Sleep for 60 seconds
	time.Sleep(60 * time.Second)

	// Update the status to "approved" (or "rejected" based on your logic)
	tx := initializers.DB.Begin()
	if err := tx.Model(&models.Transaction{}).Where("id = ?", transactionID).Update("status", "Approved").Error; err != nil {
		tx.Rollback()
		// Handle error
		return
	}
	tx.Commit()
}
