package controller

import (
	"context"
	"fmt"

	"github.com/cerbos/cerbos-sdk-go/cerbos"
	"github.com/gofiber/fiber/v2"
	"github.com/josephakayesi/go-cerbos-abac/application/dto"
	vo "github.com/josephakayesi/go-cerbos-abac/domain/core/valueobject"
	"github.com/josephakayesi/go-cerbos-abac/domain/usecase"
	internal "github.com/josephakayesi/go-cerbos-abac/internal"
	"golang.org/x/exp/slog"
)

type OrderController struct {
	OrderUsecase usecase.OrderUsecase
	CerbosClient *cerbos.GRPCClient
	Logger       slog.Logger
}

func (oc *OrderController) GetOrder(c *fiber.Ctx) error {
	ctx, cancel := context.WithCancel(c.UserContext())
	ctx = internal.SetBrowserFingerprintInContext(ctx, c)
	logId, ctx := internal.SetLogIdInContext(ctx)
	id := c.Params("id")

	defer cancel()

	oc.Logger.Info("user attempting to get order", "order_id", id)

	user := GetLoggedInUserAccessTokenPayload(c)

	order, err := oc.OrderUsecase.FindByOrderId(ctx, id)
	if err != nil {
		oc.Logger.Error("unable to find order", "log_id", logId, "err", err)
		return c.Status(400).JSON(internal.NewErrorResponse(err.Error()))
	}

	if order == nil {
		oc.Logger.Info("order not found", "log_id", logId)
		return c.Status(400).JSON(internal.NewErrorResponse("order not found"))
	}

	if user.Role.String() == vo.UserRole.String() && order.UserID != user.ID {
		oc.Logger.Info("user not authorized to get order", "log_id", logId)
		return c.Status(400).JSON(internal.NewErrorResponse(fmt.Sprintf("user not authorized to get order with id %s", id)))
	}

	oc.Logger.Info("order found successfully", "log_id", logId)

	orderDto := dto.GetOrderDto{
		ID:          order.ID,
		Description: order.Description,
		Amount:      order.Amount,
		Status:      order.Status.String(),
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
		User: dto.GetOrderUserDto{
			ID:        order.User.ID,
			FirstName: order.User.FirstName,
		},
	}

	return c.Status(200).JSON(*internal.NewSuccessResponse("order found successfully", internal.WithData(fiber.Map{"order": orderDto})))
}

func (oc *OrderController) CreateOrder(c *fiber.Ctx) error {
	ctx, cancel := context.WithCancel(c.UserContext())
	ctx = internal.SetBrowserFingerprintInContext(ctx, c)
	logId, ctx := internal.SetLogIdInContext(ctx)

	defer cancel()

	createOrderDto := &dto.CreateOrderDto{}

	if err := c.BodyParser(&createOrderDto); err != nil {
		oc.Logger.Error("unable to parse CreateOrderDto", "error", err)
		return err
	}

	oc.Logger.Info("User attempting to create order", "description", createOrderDto.Description, "amount", createOrderDto.Amount)

	createOrderDto.UserID = GetLoggedInUserAccessTokenPayload(c).ID

	_, err := oc.OrderUsecase.Create(ctx, *createOrderDto)

	if err != nil {
		oc.Logger.Error("OrderUsercase unable to create order", "log_id", logId, "err", err)
		return c.Status(400).JSON(internal.NewErrorResponse(err.Error()))
	}

	oc.Logger.Info("Order created successfully", "log_id", logId)

	return c.Status(201).JSON(*internal.NewSuccessResponse("successfully created order"))
}

func (oc *OrderController) UpdateOrder(c *fiber.Ctx) error {
	ctx, cancel := context.WithCancel(c.UserContext())
	ctx = internal.SetBrowserFingerprintInContext(ctx, c)
	logId, ctx := internal.SetLogIdInContext(ctx)
	id := c.Params("id")

	defer cancel()

	user := GetLoggedInUserAccessTokenPayload(c)

	updateOrderDto := &dto.UpdateOrderDto{}

	if err := c.BodyParser(&updateOrderDto); err != nil {
		oc.Logger.Error("unable to parse UpdateOrderDto", "error", err)
		return err
	}

	oc.Logger.Info("user attempting to update order", "description", updateOrderDto.Description, "amount", updateOrderDto.Amount, "order_id", id)

	existingOrder, err := oc.OrderUsecase.FindByOrderId(ctx, id)
	if err != nil {
		oc.Logger.Error("unable to find order", "log_id", logId, "err", err)
		return c.Status(400).JSON(internal.NewErrorResponse(err.Error()))
	}

	if existingOrder == nil {
		oc.Logger.Info("order not found", "log_id", logId)
		return c.Status(400).JSON(internal.NewErrorResponse("order not found"))
	}

	if user.ID != existingOrder.UserID {
		oc.Logger.Info("user not authorized to update order", "log_id", logId)
		return c.Status(400).JSON(internal.NewErrorResponse("user not authorized to update order"))
	}

	_, err = oc.OrderUsecase.Update(ctx, id, *updateOrderDto)

	if err != nil {
		oc.Logger.Error("OrderUsercase unable to update order", "log_id", logId, "err", err)
		return c.Status(400).JSON(internal.NewErrorResponse(err.Error()))
	}

	oc.Logger.Info("Order updated successfully", "log_id", logId)

	return c.Status(201).JSON(*internal.NewSuccessResponse("successfully updated order"))
}
