package types

import "go.mongodb.org/mongo-driver/bson"

type UpdateHotelParams struct {
	Name     string   `json:"name"`
	Location string   `json:"location"`
	Rooms    []string `json:"rooms"`
}

func (p UpdateHotelParams) ToBSON() bson.M {
	b := bson.M{}
	if p.Name != "" {
		b["name"] = p.Name
	}
	if p.Location != "" {
		b["location"] = p.Location
	}
	if len(p.Rooms) > 0 {
		b["rooms"] = p.Rooms
	}
	return b
}

type Hotel struct {
	ID       string   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string   `bson:"name" json:"name"`
	Location string   `bson:"location" json:"location"`
	Rooms    []string `bson:"rooms" json:"rooms"`
}

type RoomType int

const (
	_ RoomType = iota
	SingleBedRoomType
	DoubleBedRoomType
	ApartmentRoomType
	VipRoomType
)

type Room struct {
	ID      string   `bson:"_id,omitempty" json:"id,omitempty"`
	Type    RoomType `bson:"type" json:"type"`
	Price   float64  `bson:"price" json:"price"`
	HotelID string   `bson:"hotelID" json:"hotelID"`
}
