package models

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID            uint      `json:"accountId" db:"id"`
	UserId        uuid.UUID `json:"userId" gorm:"foreignKey:UserRefer"`
	AccountNumber uint      `json:"accountNumber" gorm:"unique"`
	BankName      string    `json:"bankName"`
	Balance       int       `json:"balance"`
	Status        string    `json:"status" gorm:"default:'Pending'"` // Default status is 'pending'
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// MarshalJSON custom implementation for JSON serialization
func (a *Account) MarshalJSON() ([]byte, error) {
	var builder strings.Builder

	builder.WriteString(`{"id":"`)
	builder.WriteString(strconv.Itoa(int(a.ID)))

	builder.WriteString(`","user_id":"`)
	builder.WriteString(a.UserId.String())

	builder.WriteString(`","account_number":"`)
	builder.WriteString(strconv.Itoa(int(a.AccountNumber)))

	builder.WriteString(`","bank_name":"`)
	builder.WriteString(a.BankName)

	builder.WriteString(`","balance":"`)
	builder.WriteString(strconv.Itoa(a.Balance))

	builder.WriteString(`","status":"`)
	builder.WriteString(a.Status)

	builder.WriteString(`","created_at":"`)
	builder.WriteString(a.CreatedAt.String())

	builder.WriteString(`","updated_at":"`)
	builder.WriteString(a.UpdatedAt.String())

	builder.WriteString(`"}`)

	return []byte(builder.String()), nil
}

// UnmarshalJSON custom implementation for JSON deserialization
func (a *Account) UnmarshalJSON(data []byte) error {
	type Alias Account
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	}

	return json.Unmarshal(data, &aux)
}
