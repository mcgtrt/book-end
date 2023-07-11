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

	userStore = store.NewMongoUserStore(client)
	userHandler = api.NewUserHandler(userStore)
}
