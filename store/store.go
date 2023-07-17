package store

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	MongoDBURL, MongoDBNAME, MongoTestDBNAME string
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
	MongoDBURL = os.Getenv("MONGO_DB_URL")
	MongoDBNAME = os.Getenv("MONGO_DB_NAME")
	MongoTestDBNAME = os.Getenv("MONGO_TEST_DB_NAME")
}
