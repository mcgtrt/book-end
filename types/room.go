package types

import "go.mongodb.org/mongo-driver/bson"

type CreateRoomParams struct {
	Type  string  `json:"type"`
	Price float64 `json:"price"`
}

func (p CreateRoomParams) ToBSON() bson.M {
	b := bson.M{}
	if p.Type != "" {
		b["type"] = p.Type
	}
	if p.Price != 0 {
		b["price"] = p.Price
	}
	return b
}

func NewRoomFromParams(params *CreateRoomParams, id string) *Room {
	return &Room{
		Type:    params.Type,
		Price:   params.Price,
		HotelID: id,
	}
}

type Room struct {
	ID      string  `bson:"_id,omitempty" json:"id,omitempty"`
	Type    string  `bson:"type" json:"type"`
	Price   float64 `bson:"price" json:"price"`
	HotelID string  `bson:"hotelID" json:"hotelID"`
}
