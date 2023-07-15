package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mcgtrt/book-end/types"
)

func AdminAuthentication(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return errUnauthorised
	}
	if !user.Admin {
		return errUnauthorised
	}
	return c.Next()
}
