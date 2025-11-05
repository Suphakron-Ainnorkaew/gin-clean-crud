package domain

import "go-clean-api/entity"

type OrderUsecase interface {
	CreateOrder(order *entity.Order, items []entity.OrderItem) error
	GetOrderByID(id uint) (*entity.Order, error)
	GetOrdersByUserID(userID uint) ([]entity.Order, error)
	UpdateOrderStatus(orderID uint, status entity.OrderStatus, shopOwnerID uint) error
	UpdatePaymentStatus(orderID uint, status entity.PaymentStatus, userID uint) error
	GetShopOrders(shopOwnerID uint) ([]entity.Order, error)
	UpdateShopOrderStatus(orderID uint, status entity.OrderStatus, shopOwnerID uint) error
	CancelShopOrder(orderID uint, shopOwnerID uint) error
	CanViewOrder(orderID uint, userID uint, userType string) (bool, error)
	DeleteOrder(id uint) error
}

type OrderRepository interface {
	Create(order *entity.Order) error
	CreateOrderItems(items []entity.OrderItem) error
	FindByID(id uint) (*entity.Order, error)
	FindByUserID(userID uint) ([]entity.Order, error)
	FindByShopID(shopID uint) ([]entity.Order, error)
	Update(order *entity.Order) error
	Delete(id uint) error
}
