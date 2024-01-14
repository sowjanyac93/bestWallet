package initializers

import "github.com/sowjanyac93/bestWallet/modules/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Account{})
	DB.AutoMigrate(&models.Transaction{})
}
