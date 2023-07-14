package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type GenericResponse struct {
	Type string `json:"type"`
	Msg  string `json:"json"`
}

func GenericResponseUnauthorised(c *fiber.Ctx) error {
	return c.Status(http.StatusUnauthorized).JSON(GenericResponse{
		Type: "error",
		Msg:  "unauthorised",
	})
}

func GenericResponseInvalidCredentials(c *fiber.Ctx) error {
	return c.Status(http.StatusBadRequest).JSON(GenericResponse{
		Type: "error",
		Msg:  "invalid credentials",
	})
}

func GenericResponseInternalServerError(c *fiber.Ctx) error {
	return c.Status(http.StatusInternalServerError).JSON(GenericResponse{
		Type: "error",
		Msg:  "internal server error",
	})
}
