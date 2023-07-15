package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/mcgtrt/book-end/store"
	"github.com/mcgtrt/book-end/types"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userStore store.UserStore
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.userStore.GetUserByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrResourceNotFound()
		}
		return ErrBadRequest()
	}
	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return ErrBadRequest()
	}
	return c.JSON(users)
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params *types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}
	if err := params.Validate(); len(err) != 0 {
		return ErrBadRequest()
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return ErrBadRequest()
	}
	insertedUser, err := h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		return ErrBadRequest()
	}
	return c.JSON(insertedUser)
}

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	var (
		id     = c.Params("id")
		params *types.UpdateUserParams
	)
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}
	if err := h.userStore.UpdateUser(c.Context(), id, params); err != nil {
		return ErrBadRequest()
	}
	return c.JSON(map[string]string{"updated": id})
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.userStore.DeleteUser(c.Context(), id); err != nil {
		return ErrBadRequest()
	}
	return c.JSON(map[string]string{"deleted": id})
}

func newUserHandler(userStore store.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}
