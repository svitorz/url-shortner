package models

import "time"

type Link struct {
	ID        uint   `gorm:"primaryKey"`
	Slug      string `gorm:"size:64;uniqueIndex"`
	TargetURL string `gorm:"column:target_url"`
	IsActive  bool   `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// FK para User
	OwnerID uint `gorm:"not null;index"`
	Owner   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:OwnerID;references:ID"`

	// Relação 1:N com Click
	//Clicks []Click `gorm:"foreignKey:LinkID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
