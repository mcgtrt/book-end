package store

import (
	"context"

	"github.com/mcgtrt/book-end/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomStore interface {
	// TODO : FINISH INTERFACE METHODS AND CORRESPONDING MONGO ROOM STORE FUNCTIONS
	InsertRoom(context.Context, *types.Room) (*types.Room, error)

	Dropper
}

type MongoRoomStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	hotelStore HotelStore
}

func (s *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	res, err := s.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	id := res.InsertedID.(primitive.ObjectID).Hex()
	room.ID = id
	if err := s.hotelStore.InsertHotelRoom(ctx, room); err != nil {
		return nil, err
	}
	return room, nil
}

func (s *MongoRoomStore) Drop(ctx context.Context) error {
	return s.coll.Drop(ctx)
}

func NewMongoRoomStore(client *mongo.Client, dbname string, hotelStore HotelStore) *MongoRoomStore {
	return &MongoRoomStore{
		client:     client,
		coll:       client.Database(dbname).Collection("rooms"),
		hotelStore: hotelStore,
	}
}
