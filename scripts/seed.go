package main

import (
	"context"

	"github.com/mcgtrt/book-end/store"
	"github.com/mcgtrt/book-end/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ctx        = context.Background()
	client     *mongo.Client
	userStore  store.UserStore
	hotelStore store.HotelStore
	roomStore  store.RoomStore
)

func main() {
	seedUser("John", "Doe", "john@doe.com", "superstrongpassword")
	seedUser("Mark", "Spencer", "mark@spencer.com", "superstrongpassword123")
	seedUser("Sabrina", "Glevesig", "sabrina@glevesig.com", "123superstrongpassword")

	seedHotel("Balenciaga", "France")
	seedHotel("Adidas", "United States")
	seedHotel("Nike", "China")
}

func init() {
	var err error
	client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(store.DBURI))
	if err != nil {
		panic(err)
	}

	userStore = store.NewMongoUserStore(client, store.DBNAME)
	hotelStore = store.NewMongoHotelStore(client, store.DBNAME)
	roomStore = store.NewMongoRoomStore(client, store.DBNAME, hotelStore)

	userStore.Drop(ctx)
	hotelStore.Drop(ctx)
	roomStore.Drop(ctx)
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
	_, err = userStore.InsertUser(ctx, user)
	return err
}

func seedHotel(name, location string) error {
	hotel := &types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []string{},
	}
	insertedHotel, err := hotelStore.InsertHotel(ctx, hotel)
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
		_, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			return err
		}
	}

	return nil
}
