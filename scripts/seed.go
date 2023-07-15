package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mcgtrt/book-end/api"
	"github.com/mcgtrt/book-end/store"
	"github.com/mcgtrt/book-end/store/fixtures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(store.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(store.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	db := store.NewMongoStore(client, store.DBNAME)

	user := fixtures.AddUser(db, "Regular", "Folk", false)
	userToken := api.CreateTokenFromUser(user)
	admin := fixtures.AddUser(db, "Geralt", "Witcher", true)
	adminToken := api.CreateTokenFromUser(admin)
	hotel := fixtures.AddHotel(db, "Kaer Morhen", "Far Noth-East", 5, nil)
	room := fixtures.AddRoom(db, "hall", hotel.ID, 399.97)
	booking := fixtures.AddBooking(db, user.ID, room.ID, 5, time.Now().AddDate(0, 0, 1), time.Now().AddDate(0, 0, 5))
	fmt.Printf("\n\nUSER TOKEN: %s\n\nADMIN TOKEN: %s\n\nBOOKING: %v\n\n", userToken, adminToken, booking)
}
