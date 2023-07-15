package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/mcgtrt/book-end/types"
)

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var (
		app         = fiber.New()
		userHandler = newUserHandler(tdb.db.User)
	)

	app.Post("/", userHandler.HandlePostUser)
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
	req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
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
