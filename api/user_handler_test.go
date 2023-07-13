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
	db *store.Store
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.db.User.Drop(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := getFiberApp(tdb.db)
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
	store := store.NewMongoStore(client, store.TestDBNAME)
	return &testdb{
		db: store,
	}
}

func getFiberApp(store *store.Store) *fiber.App {
	handler := NewHandler(store)
	app := fiber.New()
	app.Post("/auth", handler.Auth.HandleAuth)
	app.Post("/user", handler.User.HandlePostUser)
	return app
}
