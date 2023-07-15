package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mcgtrt/book-end/types"
)

func AdminAuthentication(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return ErrUnauthorised()
	}
	if !user.Admin {
		return ErrUnauthorised()
	}
	return c.Next()
}
