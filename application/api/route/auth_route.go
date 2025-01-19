package route

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/josephakayesi/go-cerbos-abac/application/api/controller"
	"github.com/josephakayesi/go-cerbos-abac/domain/repository"
	"github.com/josephakayesi/go-cerbos-abac/domain/usecase"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
)

func NewAuthRouter(timeout time.Duration, db *gorm.DB, group fiber.Router) {
	ur := repository.NewUserRespository(db)
	ac := &controller.AuthController{
		AuthUsecase: usecase.NewAuthUsecase(ur, timeout),
		Logger:      *slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("component", "auth"),
	}

	group.Post("/auth/register", ac.Register)
	group.Post("/auth/login", ac.Login)
	group.Post("/auth/refresh-token", ac.RefreshToken)
	group.Get("/auth/public-keys", ac.GetVerificationPublicKeys)
}

// func NewAuthRouter(timeout time.Duration, nq *config.NatsQueue, group fiber.Router) {
// 	// ur := repository.NewUserRespository(db)
// 	ac := &controller.AuthController{
// 		Queue:       nq,
// 		AuthUsecase: usecase.NewAuthUsecase(timeout),
// 		// AuthUsecase: usecase.NewAuthUsecase(ur, timeout),
// 		Logger: *slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("component", "auth"),
// 	}

// 	group.Post("/auth/email", ac.EmailUser)

// 	group.Post("/auth/register", ac.Register)
// 	group.Post("/auth/login", ac.Login)
// 	group.Post("/auth/refresh-token", ac.RefreshToken)
// }
