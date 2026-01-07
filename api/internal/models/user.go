package models

import (
	"time"
)

// gorm.Model definition
type User struct {
	ID        uint    `json:"id,omitempty" gorm:"primaryKey"`
	Name      string  `json:"full_name"`
	Email     *string `json:"email"`
	Password  string
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Link      []Link
}
