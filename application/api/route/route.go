package route

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type RouteOptions struct {
	Timeout time.Duration
	DB      *gorm.DB
	Engine  *fiber.App
}

func Setup(r *RouteOptions) {
	v1 := r.Engine.Group("/api/v1")

	NewHealthRouter(v1)
	NewAuthRouter(r.Timeout, r.DB, v1)
	NewOrdersRouter(r, v1)
}
