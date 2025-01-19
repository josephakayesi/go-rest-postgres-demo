package repository

import (
	"context"

	"github.com/josephakayesi/go-cerbos-abac/application/dto"
	entity "github.com/josephakayesi/go-cerbos-abac/domain/core/entity"
	"gorm.io/gorm"
)

type OrderRepository interface {
	FindByOrderId(c context.Context, id string) (*entity.Order, error)
	Create(c context.Context, order entity.Order) (*entity.Order, error)
	Update(c context.Context, id string, dto dto.UpdateOrderDto) (*entity.Order, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRespository(db *gorm.DB) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (u *orderRepository) FindByOrderId(c context.Context, id string) (*entity.Order, error) {
	order := entity.Order{}

	tx := u.db.Preload("User", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("id, first_name")
	}).Where("id = ?", id).First(&order)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return &order, nil
}

func (u *orderRepository) Create(c context.Context, order entity.Order) (*entity.Order, error) {
	tx := u.db.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return nil, tx.Error
	}

	tx.Create(&order)
	if tx.Error != nil {
		tx.Rollback()
		return nil, tx.Error
	}

	tx.Commit()
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &order, nil
}

func (u *orderRepository) Update(
	c context.Context,
	id string,
	dto dto.UpdateOrderDto,
) (*entity.Order, error) {

	order := entity.Order{}

	u.db.Model(order).Where("id = ?", id).Updates(dto)

	return &order, nil
}
