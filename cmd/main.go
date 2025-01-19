package main

import (
	// "encoding/json"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"github.com/gofiber/helmet/v2"
	route "github.com/josephakayesi/go-cerbos-abac/application/api/route"
	"github.com/josephakayesi/go-cerbos-abac/application/validation"
	entity "github.com/josephakayesi/go-cerbos-abac/domain/core/entity"
	"github.com/josephakayesi/go-cerbos-abac/infra/config"
	database "github.com/josephakayesi/go-cerbos-abac/infra/database/postgres"

	slog "golang.org/x/exp/slog"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	c := config.GetConfig()

	app := fiber.New()

	validation.NewValidator()

	app.Use(helmet.New())

	pg, err := database.NewPostgres(c)
	if err != nil {
		log.Error("unable to establish database connection", "error", err)
		panic(err)
	}

	log.Info("database connected successfully")

	err = database.CreateConnectionPool(pg)
	if err != nil {
		log.Error("unable to create connection pool", "error", err)
		panic(err)
	}

	err = database.RunMigrations(pg, &entity.User{}, &entity.Order{})
	if err != nil {
		log.Error("failed to run migrations", "error", err)
		panic(err)
	}

	app.Use(requestid.New())

	app.Use(logger.New(logger.Config{
		Format:     "${time} ${ip} ${locals:requestid} ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: time.RFC3339Nano,
	}))

	timeout := time.Duration(time.Second*1) * time.Second

	r := &route.RouteOptions{
		Timeout: timeout,
		DB:      pg,
		Engine:  app,
	}

	route.Setup(r)

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-shutdownChan
		log.Info("shutting down gracefully...")

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := app.ShutdownWithContext(ctx); err != nil {
			log.Error("server shutdown error:", err)
		}
	}()

	log.Info(fmt.Sprintf("server up and runing on port %d", c.PORT))

	err = app.Listen(fmt.Sprintf(":%d", c.PORT))
	if err != nil {
		panic(fmt.Sprintf("server was unable to start and listen on port %d", c.PORT))
	}
}
