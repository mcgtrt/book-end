package store

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DBURI      = "mongodb://localhost:27017"
	DBNAME     = "book-end"
	TestDBNAME = "book-end-test"
)

type Store struct {
	User  UserStore
	Hotel HotelStore
	Room  RoomStore
}

func NewMongoStore(client *mongo.Client, dbname string) *Store {
	user := newMongoUserStore(client, dbname)
	hotel := newMongoHotelStore(client, dbname)
	room := newMongoRoomStore(client, dbname, hotel)
	return &Store{
		User:  user,
		Hotel: hotel,
		Room:  room,
	}
}

type Dropper interface {
	Drop(context.Context) error
}
