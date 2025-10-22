package domain

import (
	"time"

	"gorm.io/gorm"
)

type Courier struct {
	ID            int            `json:"id" gorm:"primaryKey"`
	Brand         string         `json:"brand" gorm:"not null"`
	Employer_name string         `json:"employer_name" gorm:"not null"`
	Phone         string         `json:"phone" gorm:"not null"`
	Shipping_cost int            `json:"shipping_cost" gorm:"not null"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}
