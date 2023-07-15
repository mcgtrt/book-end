package api

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mcgtrt/book-end/api/middleware"
	"github.com/mcgtrt/book-end/store/fixtures"
	"github.com/mcgtrt/book-end/types"
)

func TestAdminGetBooking(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var (
		user  = fixtures.AddUser(tdb.db, "Test", "User", false)
		hotel = fixtures.AddHotel(tdb.db, "Great Hotel", "Test Location", 5, nil)
		room  = fixtures.AddRoom(tdb.db, "king size", hotel.ID, 6.99)

		fromDate = time.Now().AddDate(0, 0, 1)
		toDate   = time.Now().AddDate(0, 0, 10)
		booking  = fixtures.AddBooking(tdb.db, user.ID, room.ID, 2, fromDate, toDate)

		app            = fiber.New()
		admin          = app.Group("/", middleware.JWTAuthenticate(tdb.db.User))
		bookingHandler = newBookingHandler(tdb.db.Booking)
		req            = httptest.NewRequest("GET", "/", nil)
	)

	_ = booking
	admin.Get("/", bookingHandler.HandleGetBookings)

	res, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected to find status %d but found %d", http.StatusOK, res.StatusCode)
	}

	var bookings []*types.Booking
	if err := json.NewDecoder(res.Body).Decode(&bookings); err != nil {
		log.Fatal()
	}
}
