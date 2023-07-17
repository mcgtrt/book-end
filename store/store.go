package store

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	DBURI, DBNAME, TestDBNAME string
)

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}

func NewMongoStore(client *mongo.Client, dbname string) *Store {
	user := newMongoUserStore(client, dbname)
	hotel := newMongoHotelStore(client, dbname)
	room := newMongoRoomStore(client, dbname, hotel)
	booking := NewMongoBookingStore(client, dbname)
	return &Store{
		User:    user,
		Hotel:   hotel,
		Room:    room,
		Booking: booking,
	}
}

type Dropper interface {
	Drop(context.Context) error
}

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	DBURI = os.Getenv("MONGO_DB_URI")
	DBNAME = os.Getenv("MONGO_DB_NAME")
	TestDBNAME = os.Getenv("TEST_MONGO_DB_NAME")
}
