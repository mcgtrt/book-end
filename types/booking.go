package types

import "time"

type BookRoomParams struct {
	NumPeople int       `json:"numPeople"`
	FromDate  time.Time `json:"fromDate"`
	ToDate    time.Time `json:"toDate"`
}

func (p BookRoomParams) Validate() map[string]string {
	errors := make(map[string]string)
	now := time.Now()
	if now.After(p.FromDate) {
		errors["fromDate"] = "cannot book a room for the time in the past"
	}
	if now.After(p.ToDate) {
		errors["toDate"] = "cannot book a room to the time in the past"
	} else if p.FromDate.After(p.ToDate) {
		errors["toDate"] = "toDate must be [in time] after fromDate"
	}
	if p.NumPeople == 0 {
		errors["numPeople"] = "cannot book a room for 0 people"
	}
	return errors
}

type Booking struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    string    `bson:"userID" json:"userID"`
	RoomID    string    `bson:"roomID" json:"roomID"`
	NumPeople int       `bson:"numPeople" json:"numPeople"`
	FromDate  time.Time `bson:"fromDate" json:"fromDate"`
	ToDate    time.Time `bson:"toDate" json:"toDate"`
}
