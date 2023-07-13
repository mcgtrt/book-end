package store

import (
	"context"

	"github.com/mcgtrt/book-end/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingStore interface {
	GetBookingByID(context.Context, string) (*types.Booking, error)
	GetBookings(context.Context, bson.M) ([]*types.Booking, error)
	InsertBooking(context.Context, *types.Booking) (*types.Booking, error)

	Dropper
}

type MongoBookingStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func (s *MongoBookingStore) GetBookingByID(ctx context.Context, id string) (*types.Booking, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": oid}
	var booking *types.Booking
	if err := s.coll.FindOne(ctx, filter).Decode(&booking); err != nil {
		return nil, err
	}
	return booking, nil
}

func (s *MongoBookingStore) GetBookings(ctx context.Context, filter bson.M) ([]*types.Booking, error) {
	cur, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var bookings []*types.Booking
	if err := cur.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

func (s *MongoBookingStore) InsertBooking(ctx context.Context, room *types.Booking) (*types.Booking, error) {
	res, err := s.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return room, nil
}

func (s *MongoBookingStore) Drop(ctx context.Context) error {
	return s.coll.Drop(ctx)
}

func NewMongoBookingStore(client *mongo.Client, dbname string) *MongoBookingStore {
	return &MongoBookingStore{
		client: client,
		coll:   client.Database(dbname).Collection("bookings"),
	}
}
