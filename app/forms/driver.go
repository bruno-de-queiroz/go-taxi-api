package forms

import (
	"strconv"
	"strings"
	"taxiapp.com.br/app/exceptions"
)

type DriverForm struct {
	Name     *string `form:"name" json:"name"`
	CarPlate *string `form:"carPlate" json:"carPlate"`
}

func (f *DriverForm) IsValid() (err error) {

	e := new(exceptions.ValidationException)

	if f.Name == nil || *f.Name == "" {
		e.Put("name", "is required.")
	} else if !CARPLATE_REGEX.MatchString(*f.CarPlate) {
		e.Put("carPlate", "is invalid.")
	}

	if f.CarPlate == nil || *f.CarPlate == "" {
		e.Put("carPlate", "is required.")
	}

	if e.Size() != 0 {
		return e
	}

	return nil
}

type DriverStatusForm struct {
	Latitude  *float64 `form:"latitude" json:"latitude"`
	Longitude *float64 `form:"longitude" json:"longitude"`
	Available *bool    `form:"driverAvailable" json:"driverAvailable"`
}

func (f *DriverStatusForm) IsValid() (err error) {

	e := new(exceptions.ValidationException)

	if f.Latitude == nil {
		e.Put("latitude", "is required.")
	}

	if f.Longitude == nil {
		e.Put("longitude", "is required.")
	}

	if f.Available == nil {
		e.Put("driverAvailable", "is required.")
	}

	if e.Size() != 0 {
		return e
	}

	return nil

}

type DriverInAreaForm struct {
	Sw   string `form:"sw"`
	Ne   string `form:"ne"`
	Page *int   `form:"page"`
	Max  *int   `form:"max"`
}

func (d *DriverInAreaForm) IsValid() (err error) {

	e := new(exceptions.ValidationException)

	if d.Sw == "" {

		e.Put("sw", "is required")

	} else {

		sw := strings.Split(d.Sw, ",")

		if len(sw) != 2 {

			e.Put("sw", "is invalid")

		} else {

			if _, err := strconv.ParseFloat(sw[0], 64); err != nil {
				e.Put("sw[0]", "is invalid")
			}

			if _, err := strconv.ParseFloat(sw[1], 64); err != nil {
				e.Put("sw[1]", "is invalid")
			}
		}

	}


	if d.Ne == "" {

		e.Put("ne", "is required")

	} else {

		ne := strings.Split(d.Ne, ",")

		if len(ne) != 2 {

			e.Put("ne", "is invalid")

		} else {

			if _, err := strconv.ParseFloat(ne[0], 64); err != nil {
				e.Put("ne[0]", "is invalid")
			}

			if _, err := strconv.ParseFloat(ne[1], 64); err != nil {
				e.Put("ne[1]", "is invalid")
			}
		}
	}

	if e.Size() != 0 {
		return e
	}

	return nil

}

func (d *DriverInAreaForm) GetArea() *Area {
	var (
		err error
		lat float64
		lng float64
	)
	sw := strings.Split(d.Sw, ",")
	lat, err = strconv.ParseFloat(sw[0], 64)
	if err != nil {
		panic(err)
	}
	lng, err = strconv.ParseFloat(sw[1], 64)
	if err != nil {
		panic(err)
	}

	SW := [2]float64{lng, lat}

	ne := strings.Split(d.Ne, ",")
	lat, err = strconv.ParseFloat(ne[0], 64)
	if err != nil {
		panic(err)
	}
	lng, err = strconv.ParseFloat(ne[1], 64)
	if err != nil {
		panic(err)
	}

	NE := [2]float64{lng, lat}

	return NewArea(SW, NE)
}

type Area struct {
	SW [2]float64 `json:"sw"`
	NW [2]float64 `json:"nw"`
	SE [2]float64 `json:"se"`
	NE [2]float64 `json:"ne"`
}

func (a *Area) Polygon() [4][2]float64 {
	return [4][2]float64{a.NW, a.NE, a.SE, a.SW}
}

func NewArea(sw [2]float64, ne [2]float64) *Area {
	return &Area{
		SW: sw,
		NW: [2]float64{ne[0], sw[1]},
		NE: ne,
		SE: [2]float64{sw[0], ne[1]},
	}
}
