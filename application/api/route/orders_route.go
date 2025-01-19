package route

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/josephakayesi/go-cerbos-abac/application/api/controller"
	"github.com/josephakayesi/go-cerbos-abac/application/api/middleware"
	vo "github.com/josephakayesi/go-cerbos-abac/domain/core/valueobject"
	"github.com/josephakayesi/go-cerbos-abac/domain/repository"
	"github.com/josephakayesi/go-cerbos-abac/domain/usecase"
	"golang.org/x/exp/slog"
)

func NewOrdersRouter(r *RouteOptions, group fiber.Router) {
	or := repository.NewOrderRespository(r.DB)
	oc := &controller.OrderController{
		OrderUsecase: usecase.NewOrderUsecase(or, r.Timeout),
		Logger:       *slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("component", "orders"),
	}

	group.Post("/orders", middleware.LoadAuthorizationMiddleware(vo.UserRole, vo.AdminRole, vo.SupervisorRole), oc.CreateOrder)
	group.Put("/orders/:id", middleware.LoadAuthorizationMiddleware(vo.UserRole, vo.AdminRole, vo.SupervisorRole), oc.UpdateOrder)
	group.Get("/orders/:id", middleware.LoadAuthorizationMiddleware(vo.UserRole, vo.AdminRole, vo.SupervisorRole), oc.GetOrder)

	// group.Post("/orders/:id/approve", ac.Approve)
	// group.Post("/orders/:id/decline", ac.Decline)
}
