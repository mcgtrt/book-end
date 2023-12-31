package api

import (
	"context"
	"log"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/mcgtrt/book-end/store"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
	client *mongo.Client
	db     *store.Store
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.client.Database(store.MongoTestDBNAME).Drop(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(store.MongoDBURL))
	if err != nil {
		panic(err)
	}
	store := store.NewMongoStore(client, store.MongoTestDBNAME)
	return &testdb{
		client: client,
		db:     store,
	}
}

func getApp() *fiber.App {
	var config = fiber.Config{
		ErrorHandler: ErrorHandler,
	}
	return fiber.New(config)
}
