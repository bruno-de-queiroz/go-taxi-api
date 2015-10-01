package services

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/creativelikeadog/go-taxi-api/app/forms"
	"github.com/creativelikeadog/go-taxi-api/app/models"
	"github.com/creativelikeadog/go-taxi-api/core"
)

type DriverService struct {
	database *core.MongoAdapter
}

func (s *DriverService) New(form forms.DriverForm) (err error) {
	m := new(models.Driver)

	now := bson.Now()

	m.ID = bson.NewObjectId()
	m.Available = false
	m.Name = *form.Name
	m.CarPlate = *form.CarPlate
	m.Location = [2]float64{0.0, 0.0}
	m.Created = now
	m.Updated = now

	return s.database.Action(func(collection *mgo.Collection) error {
		return collection.Insert(m)
	})
}

func (s *DriverService) UpdateStatus(id bson.ObjectId, form forms.DriverStatusForm) (err error) {

	d := make(bson.M)
	d["location"] = [2]float64{*form.Longitude, *form.Latitude}
	d["available"] = *form.Available

	return s.database.Action(func(collection *mgo.Collection) error {
		return collection.UpdateId(id, bson.M{"$set": d})
	})

}

func (s *DriverService) All(offset int, limit int) (r []*models.Driver, err error) {
	return s.find(models.DriverQuery{Offset: offset, Limit: limit})
}

func (s *DriverService) InArea(area *forms.Area, offset int, limit int) (r []*models.DriverStatusVO, err error) {
	var (
		d []*models.Driver
	)

	d, err = s.find(models.DriverQuery{Area: area, Offset: offset, Limit: limit})
	if err != nil {
		return nil, err
	}

	r = make([]*models.DriverStatusVO, 0)

	for _, v := range d {
		r = append(r, &models.DriverStatusVO{v.ID, v.Location[1], v.Location[0], v.Available})
	}

	return r, nil
}

func (s *DriverService) find(query models.DriverQuery) (r []*models.Driver, err error) {

	r = make([]*models.Driver, 0)
	q := make(bson.M)

	if query.IDs != nil {
		q["_id"] = bson.M{"$in": query.IDs}
	}

	if query.Names != nil {
		q["name"] = bson.M{"$in": query.Names}
	}

	if query.CarPlates != nil {
		q["name"] = bson.M{"$in": query.CarPlates}
	}

	if query.Area != nil {
		q["location"] = bson.M{"$geoWithin": bson.M{"$polygon": query.Area.Polygon()}}
	}

	o := query.Offset
	l := query.Limit

	r = make([]*models.Driver, 0)

	err = s.database.Action(func(collection *mgo.Collection) error {

		qq := collection.Find(q).Skip(o)

		if l > -1 {
			qq.Limit(l)
		}

		return qq.All(&r)
	})

	if err != nil {
		return nil, err
	}

	return r, nil

}

func (s *DriverService) One(id bson.ObjectId) (m *models.Driver, err error) {

	m = new(models.Driver)

	err = s.database.Action(func(collection *mgo.Collection) error {
		return collection.FindId(id).One(m)
	})

	if err != nil {
		return nil, err
	}

	return m, nil
}

func NewDriverService(app *core.Application) *DriverService {

	config := app.Config.Database

	ad := core.NewMongoAdapter(config.Host, config.Name, "drivers")

	err := ad.Action(func(collection *mgo.Collection) error {

		if err := collection.EnsureIndex(mgo.Index{
			Key:        []string{"name", "carplate"},
			Unique:     true,
			DropDups:   true,
			Background: true,
			Sparse:     true,
		}); err != nil {
			return err
		}

		if err := collection.EnsureIndex(mgo.Index{
			Key:  []string{"$2d:point"},
			Bits: 26,
		}); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		app.Logger.Error(err)
	}

	return &DriverService{ad}
}
