package api

import "github.com/mcgtrt/book-end/store"

type Handler struct {
	Auth  *AuthHandler
	User  *UserHandler
	Hotel *HotelHandler
	Room  *RoomHandler
}

func NewHandler(store *store.Store) *Handler {
	return &Handler{
		Auth:  newAuthHandler(store.User),
		User:  newUserHandler(store.User),
		Hotel: newHotelHandler(store.Hotel, store.Room),
		Room:  newRoomHandler(store.Room),
	}
}
