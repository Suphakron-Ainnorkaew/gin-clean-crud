package entity

import (
	"time"

	"gorm.io/gorm"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type PaymentStatus string

const (
	PaymentStatusComplete PaymentStatus = "complete"
	PaymentStatusPending  PaymentStatus = "pending"
)

type Order struct {
	ID            int            `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID        int            `json:"user_id" gorm:"not null"`
	ShopID        int            `json:"shop_id" gorm:"not null"`
	CourierID     int            `json:"courier_id" gorm:"not null"`
	PaymentStatus PaymentStatus  `json:"payment_status" gorm:"type:payment_status;default:'pending';not null"`
	Status        OrderStatus    `json:"status" gorm:"type:order_status;default:'pending';not null"`
	Total         int            `json:"total" gorm:"not null"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`

	// Relations
	User       User        `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Shop       Shop        `json:"shop,omitempty" gorm:"foreignKey:ShopID"`
	Courier    Courier     `json:"courier,omitempty" gorm:"foreignKey:CourierID"`
	OrderItems []OrderItem `json:"order_items,omitempty" gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	ID        int            `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID   int            `json:"order_id" gorm:"not null"`
	ProductID int            `json:"product_id" gorm:"not null"`
	Quantity  int            `json:"quantity" gorm:"not null"`
	Price     int            `json:"price" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Relations
	Product Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
