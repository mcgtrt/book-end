package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/mcgtrt/book-end/api"
	"github.com/mcgtrt/book-end/store"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(store.DBURI))
	if err != nil {
		panic(err)
	}

	var (
		store   = store.NewMongoStore(client, store.DBNAME)
		handler = api.NewHandler(store)
		config  = fiber.Config{ErrorHandler: api.ErrorHandler}
		app     = fiber.New(config)
		apiv1   = app.Group("/api/v1", api.JWTAuthenticate(store.User))
	)

	// handle auth
	app.Post("/api/auth", handler.Auth.HandleAuth)

	// handle users
	apiv1.Get("/user/:id", handler.User.HandleGetUser)
	apiv1.Get("/user", handler.User.HandleGetUsers)
	apiv1.Post("/user", handler.User.HandlePostUser)
	apiv1.Put("/user/:id", handler.User.HandlePutUser)
	apiv1.Delete("/user/:id", handler.User.HandleDeleteUser)

	// handle hotels
	apiv1.Get("/hotel/:id", handler.Hotel.HandleGetHotel)
	apiv1.Get("/hotel", handler.Hotel.HandleGetHotels)
	apiv1.Post("/hotel", handler.Hotel.HandlePostHotel)
	apiv1.Put("/hotel/:id", handler.Hotel.HandlePutHotel)
	apiv1.Delete("/hotel/:id", handler.Hotel.HandleDeleteHotel)

	// handle hotel rooms
	apiv1.Get("/hotel/:id/room", handler.Room.HandleGetRooms)
	apiv1.Post("/hotel/:id/room", handler.Room.HandlePostRoom)

	// handle bookings
	apiv1.Get("/booking/:id", handler.Booking.HandleGetBooking)
	apiv1.Get("/booking", handler.Booking.HandleGetBookings)
	apiv1.Post("/booking/:id", handler.Booking.HandlePostBooking)
	apiv1.Put("/booking/:id/cancel", handler.Booking.HandleCancelBooking)

	listenAddr := os.Getenv("HTTP_LISTEN_ADDR")
	log.Fatal(app.Listen(listenAddr))
}
