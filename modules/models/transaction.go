package models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	UserId    int    `json:"userId" gorm:"foreignKey:UserRefer"`
	AccountId int    `json:"accountId" gorm:"foreignKey:AccountRefer"`
	Flow      string `json:"flow"`
	Amount    int    `json:"amount"`
	Details   string `json:"details"`
	Status    string `json:"status" gorm:"default:'Pending'"` // Default status is 'pending'
}
