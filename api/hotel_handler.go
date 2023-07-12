package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mcgtrt/book-end/store"
	"github.com/mcgtrt/book-end/types"
)

type HotelHandler struct {
	hotelStore store.HotelStore
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	hotel, err := h.hotelStore.GetHotelByID(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(hotel)
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	hotels, err := h.hotelStore.GetHotels(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(hotels)
}

func (h *HotelHandler) HandlePostHotel(c *fiber.Ctx) error {
	var hotel *types.Hotel
	if err := c.BodyParser(&hotel); err != nil {
		return err
	}
	insertedHotel, err := h.hotelStore.InsertHotel(c.Context(), hotel)
	if err != nil {
		return err
	}
	return c.JSON(insertedHotel)
}

func (h *HotelHandler) HandlePutHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	var params *types.UpdateHotelParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if err := h.hotelStore.UpdateHotel(c.Context(), id, params); err != nil {
		return err
	}
	return c.JSON(map[string]string{"updated": id})
}

func (h *HotelHandler) HandleDeleteHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.hotelStore.DeleteHotel(c.Context(), id); err != nil {
		return err
	}
	return c.JSON(map[string]string{"deleted": id})
}

func NewHotelHandler(hotelStore store.HotelStore) *HotelHandler {
	return &HotelHandler{
		hotelStore: hotelStore,
	}
}
