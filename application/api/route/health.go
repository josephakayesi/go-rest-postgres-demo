package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/josephakayesi/go-cerbos-abac/internal"
)

func NewHealthRouter(group fiber.Router) {
	group.Get("/health", getHealth)
	group.Get("/count", getCount)
}

func getHealth(c *fiber.Ctx) error {
	return c.SendStatus(200)
}

var count int = 0

func getCount(c *fiber.Ctx) error {

	count = count + 1

	data := struct {
		Count int `json:"count"`
	}{
		Count: count,
	}

	return c.Status(200).JSON(internal.NewSuccessResponse("okidanokh", internal.WithData(data)))
}
