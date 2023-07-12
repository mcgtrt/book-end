package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mcgtrt/book-end/api"
	"github.com/mcgtrt/book-end/store"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"err": err.Error()})
	},
}

func main() {
	listenAddr := flag.String("listenAddr", ":3000", "Listen address for hotel reservation API")
	flag.Parse()

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(store.DBURI))
	if err != nil {
		panic(err)
	}

	var (
		app         = fiber.New(config)
		apiv1       = app.Group("/api/v1")
		userStore   = store.NewMongoUserStore(client, store.DBNAME)
		userHandler = api.NewUserHandler(userStore)
	)

	// handle users
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)

	log.Fatal(app.Listen(*listenAddr))
}
