package api

import "github.com/mcgtrt/book-end/store"

type Handler struct {
	User  *UserHandler
	Hotel *HotelHandler
	Room  *RoomHandler
}

func NewHandler(store *store.Store) *Handler {
	return &Handler{
		User:  newUserHandler(store.User),
		Hotel: newHotelHandler(store.Hotel),
		Room:  newRoomHandler(store.Room),
	}
}
