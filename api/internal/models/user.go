package models

import (
	"time"
)

// gorm.Model definition
type User struct {
	ID        uint      `json:"id,omitempty" gorm:"primaryKey"`
	Name      string    `json:"full_name"`
	Email     *string   `json:"email" gorm:"uniqueIndex"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relação 1:N com Link usando OwnerID como FK no lado de Link
	Links []Link `gorm:"foreignKey:OwnerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
