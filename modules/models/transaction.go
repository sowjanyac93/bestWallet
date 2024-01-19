package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID        uint      `json:"trxId" db:"id"`
	UserId    uuid.UUID `json:"userId" gorm:"foreignKey:UserRefer"`
	AccountId uint      `json:"accountId" gorm:"foreignKey:AccountRefer"`
	Flow      string    `json:"flow"`
	Amount    int       `json:"amount"`
	Details   string    `json:"details"`
	Status    string    `json:"status" gorm:"default:'Pending'"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// MarshalJSON custom implementation for JSON serialization
func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID        uint      `json:"trxId"`
		UserId    string    `json:"userId"`
		AccountId uint      `json:"accountId"`
		Flow      string    `json:"flow"`
		Amount    int       `json:"amount"`
		Details   string    `json:"details"`
		Status    string    `json:"status"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}{
		ID:        t.ID,
		UserId:    t.UserId.String(),
		AccountId: t.AccountId,
		Flow:      t.Flow,
		Amount:    t.Amount,
		Details:   t.Details,
		Status:    t.Status,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	})
}

// UnmarshalJSON custom implementation for JSON deserialization
func (t *Transaction) UnmarshalJSON(data []byte) error {
	type Alias Transaction
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	return json.Unmarshal(data, &aux)
}
