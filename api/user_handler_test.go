package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/mcgtrt/book-end/store"
	"github.com/mcgtrt/book-end/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
	store.UserStore
}

func (d *testdb) teardown(t *testing.T) {
	if err := d.Drop(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := getFiberApp(tdb.UserStore)
	params := &types.CreateUserParams{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@doe.com",
		Password:  "strongpassword",
	}
	body, err := json.Marshal(params)
	if err != nil {
		t.Error(err)
	}
	req := httptest.NewRequest("POST", "/user", bytes.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	var user *types.User
	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		t.Error(err)
	}

	if user.ID == "" {
		t.Error("expected user ID but none found")
	}
	if params.FirstName != user.FirstName {
		t.Errorf("expected first name %s but found %s", params.FirstName, user.FirstName)
	}
	if params.LastName != user.LastName {
		t.Errorf("expected last name %s but found %s", params.LastName, user.LastName)
	}
	if params.Email != user.Email {
		t.Errorf("expected email %s but found %s", params.Email, user.Email)
	}
	if user.EncryptedPassword != "" {
		t.Error("encrypted password expected to not be included")
	}
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(store.DBURI))
	if err != nil {
		panic(err)
	}
	userStore := store.NewMongoUserStore(client, store.TestDBNAME)
	return &testdb{
		UserStore: userStore,
	}
}

func getFiberApp(userStore store.UserStore) *fiber.App {
	userHandler := NewUserHandler(userStore)
	app := fiber.New()
	app.Get("/user/:id", userHandler.HandleGetUser)
	app.Get("/user", userHandler.HandleGetUsers)
	app.Post("/user", userHandler.HandlePostUser)
	app.Put("/user/:id", userHandler.HandlePutUser)
	app.Delete("/user/:id", userHandler.HandleDeleteUser)
	return app
}
