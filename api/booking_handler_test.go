package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mcgtrt/book-end/store/fixtures"
	"github.com/mcgtrt/book-end/types"
)

func TestUserGetBooking(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var (
		nonAuthorisedUser = fixtures.AddUser(tdb.db, "tt", "uu", false)
		user              = fixtures.AddUser(tdb.db, "Test", "User", false)
		hotel             = fixtures.AddHotel(tdb.db, "Great Hotel", "Test Location", 5, nil)
		room              = fixtures.AddRoom(tdb.db, "king size", hotel.ID, 6.99)

		fromDate = time.Now().AddDate(0, 0, 1)
		toDate   = time.Now().AddDate(0, 0, 10)
		booking  = fixtures.AddBooking(tdb.db, user.ID, room.ID, 2, fromDate, toDate)

		app            = fiber.New()
		appGroup       = app.Group("/", JWTAuthenticate(tdb.db.User))
		bookingHandler = newBookingHandler(tdb.db.Booking)
	)

	appGroup.Get("/:id", bookingHandler.HandleGetBooking)
	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))

	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected to find status code %d but found %d", http.StatusOK, resp.StatusCode)
	}
	var have *types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&have); err != nil {
		t.Error(err)
	}

	// TODO: reflect.DeepEqual is always false because time stamps inserted
	// and from response are in different format.
	if have.ID != booking.ID {
		t.Errorf("expected to find booking ID %s but found %s", booking.ID, have.ID)
	}
	if have.UserID != booking.UserID {
		t.Errorf("expected to find user ID %s but found %s", booking.UserID, have.UserID)
	}
	if have.RoomID != booking.RoomID {
		t.Errorf("expected to find room ID %s but found %s", booking.RoomID, have.RoomID)
	}
	if have.NumPeople != booking.NumPeople {
		t.Errorf("expected to find number of people %d but found %d", booking.NumPeople, have.NumPeople)
	}
	if have.Canceled != booking.Canceled {
		t.Errorf("expected to find canceled %v but found %v", booking.Canceled, have.Canceled)
	}

	req = httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(nonAuthorisedUser))
	resp, err = app.Test(req)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode == http.StatusOK {
		t.Errorf("expected to not find status code %d but found %d", http.StatusOK, resp.StatusCode)
	}
}

func TestAdminGetBookings(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var (
		adminUser = fixtures.AddUser(tdb.db, "super", "admin", true)
		user      = fixtures.AddUser(tdb.db, "Test", "User", false)
		hotel     = fixtures.AddHotel(tdb.db, "Great Hotel", "Test Location", 5, nil)
		room      = fixtures.AddRoom(tdb.db, "king size", hotel.ID, 6.99)

		fromDate = time.Now().AddDate(0, 0, 1)
		toDate   = time.Now().AddDate(0, 0, 10)
		booking  = fixtures.AddBooking(tdb.db, user.ID, room.ID, 2, fromDate, toDate)

		app            = fiber.New()
		admin          = app.Group("/", JWTAuthenticate(tdb.db.User), AdminAuthentication)
		bookingHandler = newBookingHandler(tdb.db.Booking)
	)

	admin.Get("/", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(adminUser))

	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected to find %d status code but found %d", http.StatusOK, resp.StatusCode)
	}

	var bookings []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal("expected to find booking in response but an error found")
	}
	if len(bookings) != 1 {
		t.Errorf("expected to find 1 booking but %d found", len(bookings))
	}

	// TODO: reflect.DeepEqual is always false because time stamps inserted
	// and from response are in different format.
	have := bookings[0]
	if have.ID != booking.ID {
		t.Errorf("expected to find booking ID %s but found %s", booking.ID, have.ID)
	}
	if have.UserID != booking.UserID {
		t.Errorf("expected to find user ID %s but found %s", booking.UserID, have.UserID)
	}
	if have.RoomID != booking.RoomID {
		t.Errorf("expected to find room ID %s but found %s", booking.RoomID, have.RoomID)
	}
	if have.NumPeople != booking.NumPeople {
		t.Errorf("expected to find number of people %d but found %d", booking.NumPeople, have.NumPeople)
	}
	if have.Canceled != booking.Canceled {
		t.Errorf("expected to find canceled %v but found %v", booking.Canceled, have.Canceled)
	}

	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))

	res, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode == http.StatusOK {
		t.Errorf("expected to find user not authorised with status code %d but found %d", http.StatusOK, res.StatusCode)
	}
}
