package models

import "time"

type Link struct {
	ID        uint   `gorm:"primaryKey"`
	Slug      string `gorm:"size:64;uniqueIndex"`
	TargetURL string `gorm:"column:target_url"`
	IsActive  bool   `gorm:"default:true"`
	// Clicks: total de redirects já descarregado do Redis (ver package clicks).
	Clicks    int64 `gorm:"not null;default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// FK para User
	OwnerID uint `gorm:"not null;index"`
	Owner   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:OwnerID;references:ID"`
}
