package repository

import (
	"go-clean-api/domain"
	"go-clean-api/entity"

	"gorm.io/gorm"
)

type postgresOrderRepository struct {
	db *gorm.DB
}

func NewPostgresOrderRepository(db *gorm.DB) domain.OrderRepository {
	return &postgresOrderRepository{db: db}
}

func (r *postgresOrderRepository) Create(order *entity.Order) error {
	return r.db.Create(order).Error
}

func (r *postgresOrderRepository) CreateOrderItems(items []entity.OrderItem) error {
	return r.db.Create(&items).Error
}

func (r *postgresOrderRepository) FindByID(id uint) (*entity.Order, error) {
	var order entity.Order
	err := r.db.Preload("User").Preload("Shop").Preload("Courier").
		Preload("OrderItems.Product").First(&order, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

func (r *postgresOrderRepository) FindByUserID(userID uint) ([]entity.Order, error) {
	var orders []entity.Order
	err := r.db.Where("user_id = ?", userID).
		Preload("User").Preload("Shop").Preload("Courier").
		Preload("OrderItems.Product").
		Order("created_at DESC").Find(&orders).Error
	return orders, err
}

func (r *postgresOrderRepository) FindByShopID(shopID uint) ([]entity.Order, error) {
	var orders []entity.Order
	err := r.db.Where("shop_id = ?", shopID).
		Preload("User").Preload("Shop").Preload("Courier").
		Preload("OrderItems.Product").
		Order("created_at DESC").Find(&orders).Error
	return orders, err
}

func (r *postgresOrderRepository) Update(order *entity.Order) error {
	return r.db.Save(order).Error
}

func (r *postgresOrderRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Order{}, id).Error
}
