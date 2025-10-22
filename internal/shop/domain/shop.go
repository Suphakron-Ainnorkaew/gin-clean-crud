package domain

import (
	"time"

	"gorm.io/gorm"
)

type Shop struct {
	ID             int            `json:"id" gorm:"primaryKey"`
	Name           string         `json:"name" gorm:"not null"`
	Province       string         `json:"province" gorm:"not null"`
	District       string         `json:"district" gorm:"not null"`
	Subdistrict    string         `json:"subdistrict" gorm:"not null"`
	Zip_code       string         `json:"zip_code" gorm:"not null"`
	Detail_address string         `json:"detail_address" gorm:"not null"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}
