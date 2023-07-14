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

type BookingHandler struct {
	store store.BookingStore
}

func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.GetBookingByID(c.Context(), id)
	if err != nil {
		return err
	}
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(genericResponse{
			Type: "error",
			Msg:  "internal server error",
		})
	}
	if !user.Admin && user.ID != booking.UserID {
		return genericResponseUnauthorised(c)
	}
	return c.JSON(booking)
}

func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return genericResponseUnauthorised(c)
	}
	if !user.Admin {
		return genericResponseUnauthorised(c)
	}
	bookings, err := h.store.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return err
	}
	return c.JSON(bookings)
}

func (h *BookingHandler) HandlePostBooking(c *fiber.Ctx) error {
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
	insertedBooking, err := h.store.InsertBooking(c.Context(), booking)
	if err != nil {
		return err
	}
	return c.JSON(insertedBooking)
}

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	update := bson.M{
		"$set": bson.M{
			"canceled": true,
		},
	}
	if err := h.store.UpdateBooking(c.Context(), id, update); err != nil {
		return err
	}
	return c.JSON(map[string]string{"updated": id})
}

func (h *BookingHandler) isRoomAvailable(ctx context.Context, roomID string, params *types.BookRoomParams) error {
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
		"canceled": false,
	}
	bookings, err := h.store.GetBookings(ctx, where)
	if err != nil {
		return err
	}
	if len(bookings) > 0 {
		return fmt.Errorf("room already booked for selected time frame")
	}
	return nil
}

func genericResponseUnauthorised(c *fiber.Ctx) error {
	return c.Status(http.StatusUnauthorized).JSON(genericResponse{
		Type: "error",
		Msg:  "unauthorised",
	})
}

func newBookingHandler(store store.BookingStore) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}
