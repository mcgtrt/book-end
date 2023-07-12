package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mcgtrt/book-end/store"
	"github.com/mcgtrt/book-end/types"
)

type RoomHandler struct {
	roomStore store.RoomStore
}

// TODO : FINISH OTHER HANDLER METHODS WHEN FINISHED ROOM STORE

func (h *RoomHandler) HandlePostRoom(c *fiber.Ctx) error {
	var params *types.CreateRoomParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	id := c.Params("id")
	room := types.NewRoomFromParams(params, id)

	insertedRoom, err := h.roomStore.InsertRoom(c.Context(), room)
	if err != nil {
		return err
	}
	return c.JSON(insertedRoom)
}

func newRoomHandler(roomStore store.RoomStore) *RoomHandler {
	return &RoomHandler{
		roomStore: roomStore,
	}
}
