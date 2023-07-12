package main

import (
	"context"

	"github.com/mcgtrt/book-end/api"
	"github.com/mcgtrt/book-end/store"
	"github.com/mcgtrt/book-end/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ctx     = context.Background()
	client  *mongo.Client
	db      *store.Store
	handler *api.Handler
)

func main() {
	seedUser("John", "Doe", "john@doe.com", "superstrongpassword")
	seedUser("Mark", "Spencer", "mark@spencer.com", "superstrongpassword123")
	seedUser("Sabrina", "Glevesig", "sabrina@glevesig.com", "123superstrongpassword")

	seedHotel("Balenciaga", "France", 3)
	seedHotel("Adidas", "United States", 5)
	seedHotel("Nike", "China", 4)
}

func init() {
	var err error
	client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(store.DBURI))
	if err != nil {
		panic(err)
	}

	db = store.NewMongoStore(client, store.DBNAME)
	handler = api.NewHandler(db)

	db.User.Drop(ctx)
	db.Hotel.Drop(ctx)
	db.Room.Drop(ctx)
}

func seedUser(fname, lname, email, pass string) error {
	params := &types.CreateUserParams{
		FirstName: fname,
		LastName:  lname,
		Email:     email,
		Password:  pass,
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}
	_, err = db.User.InsertUser(ctx, user)
	return err
}

func seedHotel(name, location string, rating int) error {
	hotel := &types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []string{},
		Rating:   rating,
	}
	insertedHotel, err := db.Hotel.InsertHotel(ctx, hotel)
	if err != nil {
		return err
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
		_, err := db.Room.InsertRoom(ctx, &room)
		if err != nil {
			return err
		}
	}

	return nil
}
