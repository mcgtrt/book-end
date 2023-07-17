package store

import (
	"context"

	"github.com/mcgtrt/book-end/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HotelStore interface {
	GetHotelByID(context.Context, string) (*types.Hotel, error)
	GetHotels(context.Context, *options.FindOptions) ([]*types.Hotel, error)
	InsertHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	InsertHotelRoom(context.Context, *types.Room) error
	UpdateHotel(context.Context, string, *types.UpdateHotelParams) error
	DeleteHotel(context.Context, string) error

	Dropper
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func (s *MongoHotelStore) GetHotelByID(ctx context.Context, id string) (*types.Hotel, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": oid}
	res := s.coll.FindOne(ctx, filter)
	var hotel *types.Hotel
	if err := res.Decode(&hotel); err != nil {
		return nil, err
	}
	return hotel, err
}

func (s *MongoHotelStore) GetHotels(ctx context.Context, opts *options.FindOptions) ([]*types.Hotel, error) {
	cur, err := s.coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	var hotels []*types.Hotel
	if err := cur.All(ctx, &hotels); err != nil {
		return nil, err
	}
	return hotels, err
}

func (s *MongoHotelStore) InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	res, err := s.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return hotel, nil
}

func (s *MongoHotelStore) InsertHotelRoom(ctx context.Context, room *types.Room) error {
	oid, err := primitive.ObjectIDFromHex(room.HotelID)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": oid}
	update := bson.M{"$push": bson.M{"rooms": room.ID}}
	_, err = s.coll.UpdateOne(ctx, filter, update)
	return err
}

func (s *MongoHotelStore) UpdateHotel(ctx context.Context, id string, params *types.UpdateHotelParams) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": oid}
	update := bson.M{"$set": params.ToBSON()}
	_, err = s.coll.UpdateOne(ctx, filter, update)
	return err
}

func (s *MongoHotelStore) DeleteHotel(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": oid}
	_, err = s.coll.DeleteOne(ctx, filter)
	return err
}

func (s *MongoHotelStore) Drop(ctx context.Context) error {
	return s.coll.Drop(ctx)
}

func newMongoHotelStore(client *mongo.Client, dbname string) *MongoHotelStore {
	return &MongoHotelStore{
		client: client,
		coll:   client.Database(dbname).Collection("hotels"),
	}
}
