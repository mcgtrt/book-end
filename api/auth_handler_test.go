package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/mcgtrt/book-end/store/fixtures"
)

func TestAuthenticateSuccess(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var (
		insertedUser = fixtures.AddUser(tdb.db, "Test", "User", false)
		app          = getApp()
		authHandler  = newAuthHandler(tdb.db.User)
	)

	app.Post("/", authHandler.HandleAuth)
	params := &AuthParams{
		Email:    "test@user.com",
		Password: "test_user",
	}
	body, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected to find status code %d but found %d", http.StatusOK, res.StatusCode)
	}
	var authResponse AuthResponse
	if err := json.NewDecoder(res.Body).Decode(&authResponse); err != nil {
		t.Error(err)
	}
	if authResponse.Token == "" {
		t.Fatal("expected to find the JWT token but none found")
	}
	insertedUser.EncryptedPassword = ""
	if !reflect.DeepEqual(insertedUser, authResponse.User) {
		t.Fatal("expected to find inserted user but found different")
	}
}

func TestAuthenticateWrongPassword(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	fixtures.AddUser(tdb.db, "Test", "User", false)

	var (
		app         = getApp()
		authHandler = newAuthHandler(tdb.db.User)
	)

	app.Post("/", authHandler.HandleAuth)
	params := &AuthParams{
		Email:    "john@doe.com",
		Password: "superstrongpassword123",
	}
	body, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected to find status code %d but found %d", http.StatusBadRequest, res.StatusCode)
	}
	var resp Error
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatal("expected to find generic error response but none found")
	}
	invalid := ErrInvalidCredentials()
	if resp.Code != invalid.Code {
		t.Fatalf("expected to find status code %d but found %d", invalid.Code, resp.Code)
	}
	if resp.Msg != invalid.Msg {
		t.Fatalf("expected to find %s error message but found %s", invalid.Msg, resp.Msg)
	}
}
