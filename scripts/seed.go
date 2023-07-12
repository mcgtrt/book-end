package main

import (
	"context"

	"github.com/mcgtrt/book-end/api"
	"github.com/mcgtrt/book-end/store"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client      *mongo.Client
	userStore   store.UserStore
	userHandler *api.UserHandler
)

func main() {

}

func init() {
	var err error
	client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(store.DBURI))
	if err != nil {
		panic(err)
	}

	userStore = store.NewMongoUserStore(client, store.DBNAME)
	userHandler = api.NewUserHandler(userStore)
}

func seedHotel(name, location string) {
	// hotel := types.Hotel{
	// 	Name: name,
	// 	Location: location,
	// }
	// rooms := []types.Room{
	// 	types.Room{
	// 		Type: types.DoubleBedRoomType,
	// 		Price: 99.97,
	// 		HotelID:
	// 	},
	// }
}
