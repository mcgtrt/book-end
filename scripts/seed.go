package main

import (
	"context"
	"time"

	"github.com/mcgtrt/book-end/store"
	"github.com/mcgtrt/book-end/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ctx    = context.Background()
	client *mongo.Client
	db     *store.Store
)

func main() {
	user, _ := seedUser("John", "Doe", "john@doe.com", "superstrongpassword", false)
	seedUser("Mark", "Spencer", "mark@spencer.com", "superstrongpassword123", false)
	seedUser("Sabrina", "Glevesig", "sabrina@glevesig.com", "123superstrongpassword", true)

	hotel, _ := seedHotel("Adidas", "United States", 5)
	seedHotel("Puma", "France", 3)
	seedHotel("Nike", "China", 4)

	seedBooking(user.ID, hotel.Rooms[0], 5, time.Now().AddDate(0, 1, 0), time.Now().AddDate(0, 1, 5), false)
	seedBooking(user.ID, hotel.Rooms[0], 4, time.Now().AddDate(0, 1, 10), time.Now().AddDate(0, 1, 15), false)
	seedBooking(user.ID, hotel.Rooms[1], 4, time.Now().AddDate(0, 0, 20), time.Now().AddDate(0, 1, 0), false)
}

func init() {
	var err error
	client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(store.DBURI))
	if err != nil {
		panic(err)
	}

	db = store.NewMongoStore(client, store.DBNAME)

	db.User.Drop(ctx)
	db.Hotel.Drop(ctx)
	db.Room.Drop(ctx)
	db.Booking.Drop(ctx)
}

func seedUser(fname, lname, email, pass string, isAdmin bool) (*types.User, error) {
	params := &types.CreateUserParams{
		FirstName: fname,
		LastName:  lname,
		Email:     email,
		Password:  pass,
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return nil, err
	}
	user.Admin = isAdmin
	insertedUser, err := db.User.InsertUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return insertedUser, nil
}

func seedHotel(name, location string, rating int) (*types.Hotel, error) {
	hotel := &types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []string{},
		Rating:   rating,
	}
	insertedHotel, err := db.Hotel.InsertHotel(ctx, hotel)
	if err != nil {
		return nil, err
	}

	rooms := []types.Room{
		{
			Type:  types.DoubleBedRoomType,
			Price: 129.97,
		},
		{
			Type:  types.ApartmentRoomType,
			Price: 199.97,
		},
		{
			Type:  types.VipRoomType,
			Price: 299.97,
		},
	}

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		insertedRoom, err := db.Room.InsertRoom(ctx, &room)
		if err != nil {
			return nil, err
		}
		hotel.Rooms = append(hotel.Rooms, insertedRoom.ID)
	}

	return hotel, nil
}

func seedBooking(userID, roomID string, numPeople int, fromDate, toDate time.Time, canceled bool) error {
	booking := &types.Booking{
		UserID:    userID,
		RoomID:    roomID,
		NumPeople: numPeople,
		FromDate:  fromDate,
		ToDate:    toDate,
		Canceled:  canceled,
	}
	_, err := db.Booking.InsertBooking(ctx, booking)
	return err
}
