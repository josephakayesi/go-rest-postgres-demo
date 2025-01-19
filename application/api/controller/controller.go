package controller

import (
	"github.com/gofiber/fiber/v2"
	vo "github.com/josephakayesi/go-cerbos-abac/domain/core/valueobject"
)

// GetLoggedInUserAccessTokenPayload : Gets the user's access token payload from the current fiber context
func GetLoggedInUserAccessTokenPayload(c *fiber.Ctx) *vo.AccessTokenPayload {
	return c.Locals("user").(*vo.AccessTokenPayload)
}
