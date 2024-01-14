package initializers

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dbEnvVariable := os.Getenv("DB")
	database, err := gorm.Open(postgres.Open(dbEnvVariable), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to DB")
	}
	// database.AutoMigrate(&models.User{})
	DB = database
}
