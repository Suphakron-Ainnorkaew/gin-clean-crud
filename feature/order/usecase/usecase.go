package usecase

import (
	"errors"
	"go-clean-api/config"
	"go-clean-api/domain"
	"go-clean-api/entity"
)

type orderUsecase struct {
	orderRepo   domain.OrderRepository
	shopRepo    domain.ShopRepository
	courierRepo domain.CourierRepository
	userRepo    domain.UserRepository
	productRepo domain.ProductRepository
	cfg         config.ToolsConfig
}

func NewOrderUsecase(
	orderRepo domain.OrderRepository,
	shopRepo domain.ShopRepository,
	courierRepo domain.CourierRepository,
	userRepo domain.UserRepository,
	productRepo domain.ProductRepository,
	cfg config.ToolsConfig,
) domain.OrderUsecase {
	return &orderUsecase{
		orderRepo:   orderRepo,
		shopRepo:    shopRepo,
		courierRepo: courierRepo,
		userRepo:    userRepo,
		productRepo: productRepo,
		cfg:         cfg,
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

		if product.ShopID != order.ShopID {
			return errors.New("product does not belong to the specified shop")
		}

		if product.Stock < items[i].Quantity {
			u.cfg.Logrus.WithError(errors.New("insufficient stock for product: " + product.Product_name)).Error("Failed to create order")
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

func (u *orderUsecase) UpdateOrderStatus(orderID uint, status entity.OrderStatus, shopOwnerID uint) error {
	order, err := u.orderRepo.FindByID(orderID)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("order not found")
	}

	shop, err := u.shopRepo.FindByID(uint(order.ShopID))
	if err != nil {
		return err
	}
	if shop == nil {
		return errors.New("shop not found")
	}
	if shop.UserID != shopOwnerID {
		return errors.New("you don't have permission to update this order status")
	}

	order.Status = status
	return u.orderRepo.Update(order)
}

func (u *orderUsecase) UpdatePaymentStatus(orderID uint, status entity.PaymentStatus, userID uint) error {
	order, err := u.orderRepo.FindByID(orderID)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("order not found")
	}

	if order.UserID != int(userID) {
		return errors.New("you don't have permission to update this payment status")
	}

	order.PaymentStatus = status
	return u.orderRepo.Update(order)
}

func (u *orderUsecase) GetShopOrders(shopOwnerID uint) ([]entity.Order, error) {
	shop, err := u.shopRepo.FindByUserID(shopOwnerID)
	if err != nil {
		return nil, err
	}
	if shop == nil {
		return nil, errors.New("shop not found for this user")
	}

	orders, err := u.orderRepo.FindByShopID(uint(shop.ID))
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (u *orderUsecase) UpdateShopOrderStatus(orderID uint, status entity.OrderStatus, shopOwnerID uint) error {
	return u.UpdateOrderStatus(orderID, status, shopOwnerID)
}

func (u *orderUsecase) CancelShopOrder(orderID uint, shopOwnerID uint) error {
	return u.UpdateOrderStatus(orderID, entity.OrderStatusCancelled, shopOwnerID)
}

func (u *orderUsecase) GetOrdersByUserID(userID uint) ([]entity.Order, error) {
	return u.orderRepo.FindByUserID(userID)
}

func (u *orderUsecase) DeleteOrder(id uint) error {
	return u.orderRepo.Delete(id)
}

func (u *orderUsecase) CanViewOrder(orderID uint, userID uint, userType string) (bool, error) {
	order, err := u.orderRepo.FindByID(orderID)
	if err != nil {
		return false, err
	}
	if order == nil {
		return false, errors.New("order not found")
	}

	userTypeEnum := entity.UserType(userType)

	// ถ้าเป็น general user ต้องเป็นผู้สั่งซื้อ
	if userTypeEnum == entity.UserTypeGeneral {
		return order.UserID == int(userID), nil
	}

	// ถ้าเป็น shop ต้องเป็นเจ้าของร้าน
	if userTypeEnum == entity.UserTypeShop {
		shop, err := u.shopRepo.FindByID(uint(order.ShopID))
		if err != nil {
			return false, err
		}
		if shop == nil {
			return false, nil
		}
		return shop.UserID == userID, nil
	}

	return false, nil
}
