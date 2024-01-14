package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string         `json:"firstName"`
	LastName  string         `json:"lastName"`
	Email     string         `json:"email" gorm:"unique"`
	Password  string         `json:"-"`
	Address   string         `json:"address"`
	DeletedAt gorm.DeletedAt `json:"-"`
	// DateOfBirth *time.Time `json:"dateOfBirth"`
}
