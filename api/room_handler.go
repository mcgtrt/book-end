package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/mcgtrt/book-end/store"
	"github.com/mcgtrt/book-end/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	id := c.Params("id")
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(genericResponse{
			Type: "error",
			Msg:  "internal server error",
		})
	}
	var params *types.BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	errors := params.Validate()
	if len(errors) > 0 {
		return c.JSON(errors)
	}
	if err := h.isRoomAvailable(c.Context(), id, params); err != nil {
		return c.Status(http.StatusBadRequest).JSON(genericResponse{
			Type: "error",
			Msg:  err.Error(),
		})
	}
	booking := &types.Booking{
		RoomID:    id,
		UserID:    user.ID,
		NumPeople: params.NumPeople,
		FromDate:  params.FromDate,
		ToDate:    params.ToDate,
	}
	insertedBooking, err := h.store.Booking.InsertRoom(c.Context(), booking)
	if err != nil {
		return err
	}
	return c.JSON(insertedBooking)
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

func (h *RoomHandler) isRoomAvailable(ctx context.Context, roomID string, params *types.BookRoomParams) error {
	oid, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return err
	}
	where := bson.M{
		"_id": oid,
		"fromDate": bson.M{
			"$gte": params.FromDate,
		},
		"toDate": bson.M{
			"$lte": params.ToDate,
		},
	}
	bookings, err := h.store.Booking.GetBookings(ctx, where)
	if err != nil {
		return err
	}
	if len(bookings) > 0 {
		return fmt.Errorf("room already booked for selected time frame")
	}
	return nil
}

func newRoomHandler(store *store.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}
