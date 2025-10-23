package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID             int            `json:"id" gorm:"primaryKey"`
	First_name     string         `json:"first_name" gorm:"not null"`
	Last_name      string         `json:"last_name" gorm:"not null"`
	Email          string         `json:"email" gorm:"not null"`
	Province       string         `json:"province" gorm:"not null"`
	District       string         `json:"district" gorm:"not null"`
	Subdistrict    string         `json:"subdistrict" gorm:"not null"`
	Zip_code       string         `json:"zip_code" gorm:"not null"`
	Detail_address string         `json:"detail_address" gorm:"not null"`
	Phone          string         `json:"phone" gorm:"not null"`
	Password       string         `json:"password" gorm:"not null"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}
