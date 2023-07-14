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
	if err := tdb.client.Database(store.TestDBNAME).Drop(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(store.DBURI))
	if err != nil {
		panic(err)
	}
	store := store.NewMongoStore(client, store.TestDBNAME)
	return &testdb{
		client: client,
		db:     store,
	}
}

func getTestFiberApp(store *store.Store) *fiber.App {
	handler := NewHandler(store)
	app := fiber.New()
	app.Post("/auth", handler.Auth.HandleAuth)
	app.Post("/user", handler.User.HandlePostUser)
	return app
}
