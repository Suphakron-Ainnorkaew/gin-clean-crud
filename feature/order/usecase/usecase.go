package usecase

import (
	"go-clean-api/domain"
	"go-clean-api/entity"

	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
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
		err = errors.Wrap(err, "[Usecase.CreateOrder]: failed to find shop by ID")
		log.Error(err)
		return err
	}
	if shop == nil {
		err = errors.New("[Usecase.CreateOrder]: shop not found")
		log.Warn(err)
		return err
	}

	courier, err := u.courierRepo.GetByID(uint(order.CourierID))
	if err != nil {
		err = errors.Wrap(err, "[Usecase.CreateOrder]: failed to find courier by ID")
		log.Error(err)
		return err
	}
	if courier == nil {
		err = errors.New("[Usecase.CreateOrder]: courier not found")
		log.Warn(err)
		return err
	}

	user, err := u.userRepo.FindByID(uint(order.UserID))
	if err != nil {
		err = errors.Wrap(err, "[Usecase.CreateOrder]: failed to find user by ID")
		log.Error(err)
		return err
	}
	if user == nil {
		err = errors.New("[Usecase.CreateOrder]: user not found")
		log.Warn(err)
		return err
	}
	if entity.UserType(user.TypeUser) != entity.UserTypeGeneral {
		err = errors.New("[Usecase.CreateOrder]: only general users can create orders")
		log.Warn(err)
		return err
	}

	total := 0
	for i := range items {
		product, err := u.productRepo.FindProductByID(uint(items[i].ProductID))
		if err != nil {
			err = errors.Wrap(err, "[Usecase.CreateOrder]: failed to find product by ID")
			log.Error(err)
			return err
		}
		if product == nil {
			err = errors.New("[Usecase.CreateOrder]: product not found")
			log.Warn(err)
			return err
		}
		if product.ShopID != order.ShopID {
			err = errors.New("[Usecase.CreateOrder]: product does not belong to specified shop")
			log.Warn(err)
			return err
		}
		if product.Stock < items[i].Quantity {
			err = errors.New("[Usecase.CreateOrder]: insufficient stock for product: " + product.Product_name)
			log.Warn(err)
			return err
		}

		itemPrice := product.Price * items[i].Quantity
		items[i].Price = itemPrice
		total += itemPrice
	}

	total += courier.Shipping_cost
	order.Total = total

	if err := u.orderRepo.Create(order); err != nil {
		err = errors.Wrap(err, "[Usecase.CreateOrder]: failed to create order")
		log.Error(err)
		return err
	}

	for i := range items {
		items[i].OrderID = order.ID
	}

	if err := u.orderRepo.CreateOrderItems(items); err != nil {
		err = errors.Wrap(err, "[Usecase.CreateOrder]: failed to create order items")
		log.Error(err)
		return err
	}

	for _, item := range items {
		if err := u.productRepo.UpdateProductStock(uint(item.ProductID), -item.Quantity); err != nil {
			err = errors.Wrapf(err, "[Usecase.CreateOrder]: failed to update stock for productID=%d", item.ProductID)
			log.Error(err)
			return err
		}
	}

	return nil
}

func (u *orderUsecase) GetOrderByID(id uint) (*entity.Order, error) {
	order, err := u.orderRepo.FindByID(id)
	if err != nil {
		err = errors.Wrap(err, "[Usecase.GetOrderByID]: failed to get order by ID")
		log.Error(err)
		return nil, err
	}
	if order == nil {
		err = errors.New("[Usecase.GetOrderByID]: order not found")
		log.Warn(err)
		return nil, err
	}
	return order, nil
}

func (u *orderUsecase) UpdateOrderStatus(orderID uint, status entity.OrderStatus, shopOwnerID uint) error {
	order, err := u.orderRepo.FindByID(orderID)
	if err != nil {
		err = errors.Wrap(err, "[Usecase.UpdateOrderStatus]: failed to find order")
		log.Error(err)
		return err
	}
	if order == nil {
		err = errors.New("[Usecase.UpdateOrderStatus]: order not found")
		log.Warn(err)
		return err
	}

	shop, err := u.shopRepo.FindByID(uint(order.ShopID))
	if err != nil {
		err = errors.Wrap(err, "[Usecase.UpdateOrderStatus]: failed to find shop")
		log.Error(err)
		return err
	}
	if shop == nil {
		err = errors.New("[Usecase.UpdateOrderStatus]: shop not found")
		log.Warn(err)
		return err
	}
	if shop.UserID != shopOwnerID {
		err = errors.New("[Usecase.UpdateOrderStatus]: unauthorized update attempt")
		log.Warn(err)
		return err
	}

	if err := u.orderRepo.Update(order); err != nil {
		err = errors.Wrap(err, "[Usecase.UpdateOrderStatus]: failed to update order status")
		log.Error(err)
		return err
	}

	return nil
}

func (u *orderUsecase) UpdatePaymentStatus(orderID uint, status entity.PaymentStatus, userID uint) error {
	order, err := u.orderRepo.FindByID(orderID)
	if err != nil {
		err = errors.Wrap(err, "[Usecase.UpdatePaymentStatus]: failed to find order")
		log.Error(err)
		return err
	}
	if order == nil {
		err = errors.New("[Usecase.UpdatePaymentStatus]: order not found")
		log.Warn(err)
		return err
	}

	if order.UserID != int(userID) {
		err = errors.New("[Usecase.UpdatePaymentStatus]: unauthorized user")
		log.Warn(err)
		return err
	}

	order.PaymentStatus = status
	if err := u.orderRepo.Update(order); err != nil {
		err = errors.Wrap(err, "[Usecase.UpdatePaymentStatus]: failed to update payment status")
		log.Error(err)
		return err
	}

	return nil
}

func (u *orderUsecase) GetShopOrders(shopOwnerID uint) ([]entity.Order, error) {
	shop, err := u.shopRepo.FindByUserID(shopOwnerID)
	if err != nil {
		err = errors.Wrap(err, "[Usecase.GetShopOrders]: failed to get shop by ownerID")
		log.Error(err)
		return nil, err
	}
	if shop == nil {
		err = errors.New("[Usecase.GetShopOrders]: shop not found for this user")
		log.Warn(err)
		return nil, err
	}

	orders, err := u.orderRepo.FindByShopID(uint(shop.ID))
	if err != nil {
		err = errors.Wrap(err, "[Usecase.GetShopOrders]: failed to get orders by shopID")
		log.Error(err)
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
	orders, err := u.orderRepo.FindByUserID(userID)
	if err != nil {
		err = errors.Wrap(err, "[Usecase.GetOrdersByUserID]: failed to get orders by userID")
		log.Error(err)
		return nil, err
	}
	return orders, nil
}

func (u *orderUsecase) DeleteOrder(id uint) error {
	if err := u.orderRepo.Delete(id); err != nil {
		err = errors.Wrap(err, "[Usecase.DeleteOrder]: failed to delete order")
		log.Error(err)
		return err
	}
	return nil
}

func (u *orderUsecase) CanViewOrder(orderID uint, userID uint, userType string) (bool, error) {
	order, err := u.orderRepo.FindByID(orderID)
	if err != nil {
		err = errors.Wrap(err, "[Usecase.CanViewOrder]: failed to find order")
		log.Error(err)
		return false, err
	}
	if order == nil {
		err = errors.New("[Usecase.CanViewOrder]: order not found")
		log.Warn(err)
		return false, err
	}

	userTypeEnum := entity.UserType(userType)

	if userTypeEnum == entity.UserTypeGeneral {
		return order.UserID == int(userID), nil
	}

	if userTypeEnum == entity.UserTypeShop {
		shop, err := u.shopRepo.FindByID(uint(order.ShopID))
		if err != nil {
			err = errors.Wrap(err, "[Usecase.CanViewOrder]: failed to get shop by ID")
			log.Error(err)
			return false, err
		}
		if shop == nil {
			err = errors.New("[Usecase.CanViewOrder]: shop not found")
			log.Warn(err)
			return false, err
		}
		return shop.UserID == userID, nil
	}

	err = errors.New("[Usecase.CanViewOrder]: unsupported user type")
	log.Warn(err)
	return false, err
}
