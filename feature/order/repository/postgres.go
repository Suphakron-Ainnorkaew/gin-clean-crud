package repository

import (
	"go-clean-api/domain"
	"go-clean-api/entity"

	"github.com/pkg/errors"

	"gorm.io/gorm"
)

type postgresOrderRepository struct {
	db *gorm.DB
}

func NewPostgresOrderRepository(db *gorm.DB) domain.OrderRepository {
	return &postgresOrderRepository{db: db}
}

func (r *postgresOrderRepository) Create(order *entity.Order) error {
	if err := r.db.Create(order).Error; err != nil {
		return errors.Wrap(err, "[OrderRepository.Create]: unable to create order")
	}
	return nil
}

func (r *postgresOrderRepository) CreateOrderItems(items []entity.OrderItem) error {
	if err := r.db.Create(&items).Error; err != nil {
		return errors.Wrap(err, "[OrderRepository.CreateOrderItems]: unable to create order items")
	}
	return nil
}

func (r *postgresOrderRepository) FindByID(id uint) (*entity.Order, error) {
	var order entity.Order
	if err := r.db.Preload("User").Preload("Shop").Preload("Courier").
		Preload("OrderItems.Product").First(&order, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errors.Wrap(err, "[OrderRepository.FindByID]: unable to find order items")
	}
	return &order, nil
}

func (r *postgresOrderRepository) FindByUserID(userID uint) ([]entity.Order, error) {
	var orders []entity.Order
	if err := r.db.Where("user_id = ?", userID).
		Preload("User").Preload("Shop").Preload("Courier").
		Preload("OrderItems.Product").
		Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, errors.Wrap(err, "[OrderRepository.FindByUserID]: unable to find order by user id")
	}
	return orders, nil
}

func (r *postgresOrderRepository) FindByShopID(shopID uint) ([]entity.Order, error) {
	var orders []entity.Order
	if err := r.db.Where("shop_id = ?", shopID).
		Preload("User").Preload("Shop").Preload("Courier").
		Preload("OrderItems.Product").
		Order("created_at DESC").Find(&orders).Error; err != nil {
		return orders, errors.Wrap(err, "[OrderRepository.FindByShopID]: unable to find order by shop id")
	}
	return orders, nil
}

func (r *postgresOrderRepository) Update(order *entity.Order) error {
	if err := r.db.Model(&entity.Order{}).Where("id = ?", order.ID).Updates(&order).Error; err != nil {
		return errors.Wrap(err, "[OrderRepository.Update]: unable to update shop")
	}
	return nil
}

func (r *postgresOrderRepository) Delete(id uint) error {
	if err := r.db.Delete(&entity.Order{}, id).Error; err != nil {
		return errors.Wrap(err, "[OrderRepository.Delete]: unable to Delete shop")
	}
	return nil
}
