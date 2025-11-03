package usecase

import (
	"errors"
	"go-clean-api/domain"
	"go-clean-api/entity"
)

type orderUsecase struct {
	orderRepo   domain.OrderRepository
	shopRepo    domain.ShopRepository
	courierRepo domain.CourierRepository
	userRepo    domain.UserRepository
	productRepo domain.ProductRepository
}

func NewOrderUsecase(
	orderRepo domain.OrderRepository,
	shopRepo domain.ShopRepository,
	courierRepo domain.CourierRepository,
	userRepo domain.UserRepository,
	productRepo domain.ProductRepository,
) domain.OrderUsecase {
	return &orderUsecase{
		orderRepo:   orderRepo,
		shopRepo:    shopRepo,
		courierRepo: courierRepo,
		userRepo:    userRepo,
		productRepo: productRepo,
	}
}

func (u *orderUsecase) CreateOrder(order *entity.Order, items []entity.OrderItem) error {
	shop, err := u.shopRepo.FindByID(uint(order.ShopID))
	if err != nil {
		return err
	}
	if shop == nil {
		return errors.New("shop not found")
	}

	courier, err := u.courierRepo.GetByID(uint(order.CourierID))
	if err != nil {
		return err
	}
	if courier == nil {
		return errors.New("courier not found")
	}

	usr, err := u.userRepo.FindByID(uint(order.UserID))
	if err != nil {
		return err
	}
	if usr == nil {
		return errors.New("user not found")
	}
	if entity.UserType(usr.TypeUser) != entity.UserTypeGeneral {
		return errors.New("only general users can create orders")
	}

	total := 0
	for i := range items {
		product, err := u.productRepo.FindProductByID(uint(items[i].ProductID))
		if err != nil {
			return err
		}
		if product == nil {
			return errors.New("product not found")
		}

		if product.Stock < items[i].Quantity {
			return errors.New("insufficient stock for product: " + product.Product_name)
		}
		itemPrice := product.Price * items[i].Quantity
		items[i].Price = itemPrice
		total += itemPrice
	}

	total += courier.Shipping_cost
	order.Total = total

	if err := u.orderRepo.Create(order); err != nil {
		return err
	}

	for i := range items {
		items[i].OrderID = order.ID
	}

	if err := u.orderRepo.CreateOrderItems(items); err != nil {
		return err
	}

	for _, item := range items {
		if err := u.productRepo.UpdateProductStock(uint(item.ProductID), -item.Quantity); err != nil {
			return err
		}
	}

	return nil
}

func (u *orderUsecase) GetOrderByID(id uint) (*entity.Order, error) {
	return u.orderRepo.FindByID(id)
}

func (u *orderUsecase) UpdateOrderStatus(orderID uint, status entity.OrderStatus) error {
	order, err := u.orderRepo.FindByID(orderID)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("order not found")
	}

	order.Status = status
	return u.orderRepo.Update(order)
}

func (u *orderUsecase) UpdatePaymentStatus(orderID uint, status entity.PaymentStatus) error {
	order, err := u.orderRepo.FindByID(orderID)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("order not found")
	}

	order.PaymentStatus = status
	return u.orderRepo.Update(order)
}

func (u *orderUsecase) GetOrdersByUserID(userID uint) ([]entity.Order, error) {
	return u.orderRepo.FindByUserID(userID)
}

func (u *orderUsecase) DeleteOrder(id uint) error {
	return u.orderRepo.Delete(id)
}
