package models

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Username    *string   `json:"username" db:"username"`
	FirstName   *string   `json:"first_name" db:"first_name"`
	LastName    *string   `json:"last_name" db:"last_name"`
	Email       string    `json:"email" db:"email"`
	PhoneNumber *string   `json:"phone_number" db:"phone_number"`
	Password    string    `json:"password" db:"password"`
	Address     string    `json:"address" db:"address"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func (u *User) MarshalJSON() ([]byte, error) {
	var builder strings.Builder

	builder.WriteString(`{"id":"`)
	builder.WriteString(u.ID.String())
	if u.Username != nil {
		builder.WriteString(`","username":"`)
		builder.WriteString(*(u.Username))
	}

	if u.FirstName != nil {
		builder.WriteString(`","first_name":"`)
		builder.WriteString(*(u.FirstName))
	}

	if u.LastName != nil {
		builder.WriteString(`","last_name":"`)
		builder.WriteString(*(u.LastName))
	}

	builder.WriteString(`","email":"`)
	builder.WriteString(u.Email)
	if u.PhoneNumber != nil {
		builder.WriteString(`","phone_number":"`)
		builder.WriteString(*(u.PhoneNumber))
	}

	builder.WriteString(`","address":"`)
	builder.WriteString(u.Address)

	builder.WriteString(`","created_at":"`)
	builder.WriteString(u.CreatedAt.String())

	builder.WriteString(`","updated_at":"`)
	builder.WriteString(u.UpdatedAt.String())

	builder.WriteString(`"}`)

	return []byte(builder.String()), nil
}

func (u *User) UnmarshalJSON(data []byte) error {
	type Alias User
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(u),
	}

	return json.Unmarshal(data, &aux)
}
