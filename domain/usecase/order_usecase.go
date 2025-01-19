package usecase

import (
	"context"
	"time"

	"github.com/josephakayesi/go-cerbos-abac/application/dto"
	"github.com/josephakayesi/go-cerbos-abac/application/validation"
	entity "github.com/josephakayesi/go-cerbos-abac/domain/core/entity"
	"github.com/josephakayesi/go-cerbos-abac/domain/repository"
	"github.com/josephakayesi/go-cerbos-abac/internal"
)

type OrderUsecase interface {
	FindByOrderId(c context.Context, id string) (*entity.Order, error)
	Create(c context.Context, l dto.CreateOrderDto) (*entity.Order, error)
	Update(c context.Context, id string, l dto.UpdateOrderDto) (*entity.Order, error)
}

type orderUsecase struct {
	orderRepository repository.OrderRepository
	contextTimeout  time.Duration
}

func NewOrderUsecase(orderRepository repository.OrderRepository, timeout time.Duration) OrderUsecase {
	return &orderUsecase{
		orderRepository: orderRepository,
		contextTimeout:  timeout,
	}
}

func (ou *orderUsecase) FindByOrderId(c context.Context, id string) (*entity.Order, error) {
	_, cancel := context.WithTimeout(c, ou.contextTimeout)
	defer cancel()

	order, err := ou.orderRepository.FindByOrderId(c, id)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (ou *orderUsecase) Create(c context.Context, co dto.CreateOrderDto) (*entity.Order, error) {
	ctx, cancel := context.WithTimeout(c, ou.contextTimeout)
	defer cancel()

	order := &entity.Order{}

	err := validation.CreateOrderDtoValidation(&co)

	if err != nil {
		return nil, err
	}

	order.ID = internal.GenerateOrderId()
	order.Description = co.Description
	order.Amount = co.Amount
	order.UserID = co.UserID

	_, err = ou.orderRepository.Create(ctx, *order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (ou *orderUsecase) Update(c context.Context, id string, dto dto.UpdateOrderDto) (*entity.Order, error) {
	ctx, cancel := context.WithTimeout(c, ou.contextTimeout)
	defer cancel()

	err := validation.UpdateOrderDtoValidation(&dto)

	if err != nil {
		return nil, err
	}

	order, err := ou.orderRepository.Update(ctx, id, dto)
	if err != nil {
		return nil, err
	}

	return order, nil
}
