package models

import (
	"time"

	"github.com/creativelikeadog/go-taxi-api/app/forms"
	"gopkg.in/mgo.v2/bson"
)

type Driver struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Name      string        `json:"name"`
	CarPlate  string        `json:"car_plate"`
	Location  [2]float64    `json:"location"`
	Available bool          `json:"available"`
	Created   time.Time     `json:"created_at"`
	Updated   time.Time     `json:"updated_at"`
}

type DriverQuery struct {
	IDs       []bson.ObjectId
	Names     []string
	CarPlates []string
	Area      *forms.Area
	Offset    int
	Limit     int
}

type DriverStatusVO struct {
	ID        bson.ObjectId `json:"driver_id"`
	Latitude  float64       `json:"latitude"`
	Longitude float64       `json:"longitude"`
	Available bool          `json:"driver_available"`
}
