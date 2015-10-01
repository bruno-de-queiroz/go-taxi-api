package core

import "gopkg.in/mgo.v2"

type MongoAdapter struct {
	host       string
	database   string
	collection string
	session    *mgo.Session
}

func (C *MongoAdapter) clone() (session *mgo.Session, err error) {
	if C.session == nil {
		C.session, err = mgo.Dial(C.host)
		if err != nil {
			return nil, err
		}

		C.session.SetMode(mgo.Monotonic, true)
	}

	return C.session.Copy(), nil
}

func (C *MongoAdapter) Action(a func(*mgo.Collection) error) (err error) {
	session, err := C.clone()
	if err != nil {
		return err
	}

	defer session.Close()
	return a(session.DB(C.database).C(C.collection))
}

func NewMongoAdapter(host string, database string, collection string) *MongoAdapter {
	return &MongoAdapter{host, database, collection, nil}
}
