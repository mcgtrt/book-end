package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mcgtrt/book-end/store"
	"github.com/mcgtrt/book-end/types"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HotelHandler struct {
	hotelStore store.HotelStore
	roomStore  store.RoomStore
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	hotel, err := h.hotelStore.GetHotelByID(c.Context(), id)
	if err != nil {
		return ErrInvalidID()
	}
	return c.JSON(hotel)
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var params types.HotelQueryParams
	if err := c.QueryParser(&params); err != nil {
		return ErrBadRequest()
	}
	hotels, err := h.hotelStore.GetHotels(c.Context(), createPaginationOptions(&params))
	if err != nil {
		return ErrBadRequest()
	}
	if params.Rooms {
		var hotelsWithRooms []*types.HotelWithRooms
		for _, hotel := range hotels {
			rooms, err := h.roomStore.GetRooms(c.Context(), hotel.ID)
			if err != nil {
				return ErrBadRequest()
			}
			hotelWithRooms := hotel.MakeHotelWithRooms(rooms)
			hotelsWithRooms = append(hotelsWithRooms, hotelWithRooms)
		}
		return c.JSON(hotelsWithRooms)
	}
	return c.JSON(hotels)
}

func (h *HotelHandler) HandlePostHotel(c *fiber.Ctx) error {
	var hotel *types.Hotel
	if err := c.BodyParser(&hotel); err != nil {
		return ErrBadRequest()
	}
	insertedHotel, err := h.hotelStore.InsertHotel(c.Context(), hotel)
	if err != nil {
		return ErrBadRequest()
	}
	return c.JSON(insertedHotel)
}

func (h *HotelHandler) HandlePutHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	var params *types.UpdateHotelParams
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}
	if err := h.hotelStore.UpdateHotel(c.Context(), id, params); err != nil {
		return ErrBadRequest()
	}
	return c.JSON(map[string]string{"updated": id})
}

func (h *HotelHandler) HandleDeleteHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.hotelStore.DeleteHotel(c.Context(), id); err != nil {
		return ErrBadRequest()
	}
	return c.JSON(map[string]string{"deleted": id})
}

func newHotelHandler(hotelStore store.HotelStore, roomStore store.RoomStore) *HotelHandler {
	return &HotelHandler{
		hotelStore: hotelStore,
		roomStore:  roomStore,
	}
}

func createPaginationOptions(params *types.HotelQueryParams) *options.FindOptions {
	opts := &options.FindOptions{}
	opts.SetSkip((params.Page - 1) * params.Limit)
	opts.SetLimit(params.Limit)
	return opts
}
