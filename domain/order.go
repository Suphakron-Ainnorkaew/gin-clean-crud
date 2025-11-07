package domain

import (
	"go-clean-api/entity"

	"github.com/sirupsen/logrus"
)

type OrderUsecase interface {
	CreateOrder(log *logrus.Entry, order *entity.Order, items []entity.OrderItem) error
	GetOrderByID(log *logrus.Entry, id uint) (*entity.Order, error)
	GetOrdersByUserID(log *logrus.Entry, userID uint) ([]entity.Order, error)
	UpdateOrderStatus(log *logrus.Entry, orderID uint, status entity.OrderStatus, shopOwnerID uint) error
	UpdatePaymentStatus(log *logrus.Entry, orderID uint, status entity.PaymentStatus, userID uint) error
	GetShopOrders(log *logrus.Entry, shopOwnerID uint) ([]entity.Order, error)
	UpdateShopOrderStatus(log *logrus.Entry, orderID uint, status entity.OrderStatus, shopOwnerID uint) error
	CancelShopOrder(log *logrus.Entry, orderID uint, shopOwnerID uint) error
	CanViewOrder(log *logrus.Entry, orderID uint, userID uint, userType string) (bool, error)
	DeleteOrder(log *logrus.Entry, id uint) error
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
