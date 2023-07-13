package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mcgtrt/book-end/store"
	"github.com/mcgtrt/book-end/types"
)

type RoomHandler struct {
	store *store.Store
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	rooms, err := h.store.Room.GetRooms(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}

func (h *RoomHandler) HandlePostRoom(c *fiber.Ctx) error {
	var params *types.CreateRoomParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	id := c.Params("id")
	room := types.NewRoomFromParams(params, id)

	insertedRoom, err := h.store.Room.InsertRoom(c.Context(), room)
	if err != nil {
		return err
	}
	return c.JSON(insertedRoom)
}

func newRoomHandler(store *store.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}
