package entity

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID           int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Product_name string `json:"product_name" gorm:"not null"`
	Price        int    `json:"price" gorm:"not null"`
	Stock        int    `json:"stock" gorm:"not null"`

	ShopID int `json:"shop_id" gorm:"not null"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
