package api

import (
	"context"
	"log"
	"net/http"
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

func getApp() *fiber.App {
	var config = fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if apiErr, ok := err.(Error); ok {
				return c.Status(apiErr.Code).JSON(apiErr)
			}
			apiErr := NewError(http.StatusInternalServerError, err.Error())
			return c.Status(apiErr.Code).JSON(apiErr)
		},
	}
	return fiber.New(config)
}
