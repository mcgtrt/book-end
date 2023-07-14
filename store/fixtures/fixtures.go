package fixtures

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/mcgtrt/book-end/store"
	"github.com/mcgtrt/book-end/types"
)

func AddBooking(db *store.Store, userID, roomID string, numPeople int, fromDate, toDate time.Time) *types.Booking {
	booking := &types.Booking{
		UserID:    userID,
		RoomID:    roomID,
		NumPeople: numPeople,
		FromDate:  fromDate,
		ToDate:    toDate,
		Canceled:  false,
	}
	insertedBooking, err := db.Booking.InsertBooking(context.Background(), booking)
	if err != nil {
		log.Fatal(err)
	}
	return insertedBooking
}

func AddHotel(db *store.Store, name, loc string, rating int, rooms []string) *types.Hotel {
	ids := rooms
	if rooms == nil {
		ids = []string{}
	}
	hotel := &types.Hotel{
		Name:     name,
		Location: loc,
		Rooms:    ids,
		Rating:   rating,
	}
	insertedHotel, err := db.Hotel.InsertHotel(context.Background(), hotel)
	if err != nil {
		log.Fatal(err)
	}
	return insertedHotel
}

func AddRoom(db *store.Store, roomType, hotelID string, price float64) *types.Room {
	room := &types.Room{
		Type:    roomType,
		Price:   129.97,
		HotelID: hotelID,
	}

	insertedRoom, err := db.Room.InsertRoom(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}

	return insertedRoom
}

func AddUser(db *store.Store, fn, ln string, admin bool) *types.User {
	fnLower := strings.ToLower(fn)
	lnLower := strings.ToLower(ln)
	params := &types.CreateUserParams{
		FirstName: fn,
		LastName:  ln,
		Email:     fmt.Sprintf("%s@%s.com", fnLower, lnLower),
		Password:  fmt.Sprintf("%s_%s", fnLower, lnLower),
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		log.Fatal(err)
	}
	user.Admin = admin
	insertedUser, err := db.User.InsertUser(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}
	return insertedUser
}
