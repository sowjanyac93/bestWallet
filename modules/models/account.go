package models

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	UserId        int    `json:"userId" gorm:"foreignKey:UserRefer"`
	AccountNumber int    `json:"accountNumber" gorm:"unique"`
	BankName      string `json:"bankName"`
	Balance       int    `json:"balance"`
	Status        string `json:"status" gorm:"default:'Pending'"` // Default status is 'pending'
}
