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
	insertedUser := fixtures.AddUser(tdb.db, "Test", "User", false)

	app := getTestFiberApp(tdb.db)
	params := &AuthParams{
		Email:    "test@user.com",
		Password: "test_user",
	}
	body, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(body))
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

	app := getTestFiberApp(tdb.db)
	params := &AuthParams{
		Email:    "john@doe.com",
		Password: "superstrongpassword123",
	}
	body, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected to find status code %d but found %d", http.StatusBadRequest, res.StatusCode)
	}
	var resp GenericResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatal("expected to find a generic response but none found")
	}
	if resp.Type != "error" {
		t.Fatalf("expected to find response type error but found %s", resp.Type)
	}
	if resp.Msg != "invalid credentials" {
		t.Fatalf("expected to find <invalid credentials> error message but found %s", resp.Msg)
	}
}
