package domain

import "go-clean-api/entity"

type OrderUsecase interface {
	CreateOrder(order *entity.Order, items []entity.OrderItem) error
	GetOrderByID(id uint) (*entity.Order, error)
	GetOrdersByUserID(userID uint) ([]entity.Order, error)
	UpdateOrderStatus(orderID uint, status entity.OrderStatus) error
	UpdatePaymentStatus(orderID uint, status entity.PaymentStatus) error
	DeleteOrder(id uint) error
}

type OrderRepository interface {
	Create(order *entity.Order) error
	CreateOrderItems(items []entity.OrderItem) error
	FindByID(id uint) (*entity.Order, error)
	FindByUserID(userID uint) ([]entity.Order, error)
	Update(order *entity.Order) error
	Delete(id uint) error
}
