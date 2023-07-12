package types

import "go.mongodb.org/mongo-driver/bson"

type HotelQueryParams struct {
	Rooms  bool
	Rating int
}

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
	Rating   int      `bson:"rating" json:"rating"`
}

func (h *Hotel) MakeHotelWithRooms(rooms []*Room) *HotelWithRooms {
	return &HotelWithRooms{
		ID:       h.ID,
		Name:     h.Name,
		Location: h.Location,
		Rooms:    rooms,
		Rating:   h.Rating,
	}
}

type HotelWithRooms struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Location string  `json:"location"`
	Rooms    []*Room `json:"rooms"`
	Rating   int     `json:"rating"`
}
