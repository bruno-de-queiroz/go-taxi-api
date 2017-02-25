package services

import (
	"errors"

	"github.com/creativelikeadog/go-taxi-api/app/forms"
	"github.com/creativelikeadog/go-taxi-api/app/mailers"
	"github.com/creativelikeadog/go-taxi-api/app/models"
	"github.com/creativelikeadog/go-taxi-api/core"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type TokenType int

const (
	ACCESS_TOKEN TokenType = iota
	RESET_TOKEN
)

var (
	tokens = [...]string{"accessToken", "resetToken"}
)

func (t TokenType) String() string {
	if len(tokens) > int(t) {
		return tokens[t]
	}
	return ""
}

type UserService struct {
	logger   *core.Logger
	database *core.MongoAdapter
	sender   *core.EmailSender
}

func (s *UserService) encryptPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (s *UserService) ByEmailAndPassword(email string, password string) (m *models.User, err error) {

	m = new(models.User)

	err = s.database.Action(func(collection *mgo.Collection) error {
		return collection.Find(bson.M{"email": email, "resetToken": bson.M{"$exists": false}}).One(m)
	})

	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(m.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (s *UserService) ByEmail(email string) (m *models.User, err error) {
	m = new(models.User)

	err = s.database.Action(func(collection *mgo.Collection) error {
		return collection.Find(bson.M{"email": email}).One(m)
	})

	if err != nil {
		return nil, err
	}

	return m, nil
}

func (s *UserService) One(id bson.ObjectId) (m *models.User, err error) {
	m = new(models.User)

	err = s.database.Action(func(collection *mgo.Collection) error {
		return collection.FindId(id).One(m)
	})

	if err != nil {
		return nil, err
	}

	return m, nil
}

func (s *UserService) NewPassword(user bson.ObjectId, password string) (err error) {

	hash, err := s.encryptPassword(password)
	if err != nil {
		return err
	}

	return s.database.Action(func(collection *mgo.Collection) error {
		return collection.UpdateId(user, bson.M{"$set": bson.M{"password": string(hash)}, "$unset": bson.M{"resetToken": ""}})
	})
}

func (s *UserService) RemoveToken(user bson.ObjectId, tokenType TokenType) (err error) {

	t := tokenType.String()

	if t == "" {
		return errors.New("Token type not allowed")
	}

	return s.database.Action(func(collection *mgo.Collection) error {
		return collection.UpdateId(user, bson.M{"$unset": bson.M{t: ""}})
	})
}

func (s *UserService) SaveToken(user bson.ObjectId, tokenType TokenType, token string) (err error) {
	t := tokenType.String()

	if t == "" {
		return errors.New("Token type not allowed")
	}

	return s.database.Action(func(collection *mgo.Collection) error {
		return collection.UpdateId(user, bson.M{"$set": bson.M{t: token}})
	})
}

func (s *UserService) New(form *forms.UserForm) (m *models.User, err error) {

	var (
		hash []byte
	)

	hash, err = s.encryptPassword(*form.Password)
	if err != nil {
		return nil, err
	}

	now := bson.Now()

	m = new(models.User)
	m.ID = bson.NewObjectId()
	m.Name = *form.Name
	m.Email = *form.Email
	m.Password = string(hash)
	m.Created = now
	m.Updated = now

	err = s.database.Action(func(collection *mgo.Collection) error {
		return collection.Insert(m)
	})

	if err != nil {
		return nil, err
	}

	err = s.sender.Send(mailers.NewUserRegisteredEmail(m))
	if err != nil {
		s.logger.Error(err)
	}

	return m, nil
}

func NewUserService(app *core.Application) *UserService {

	ad := core.NewMongoAdapter(app.Config.Database.Host, app.Config.Database.Name, "users")

	err := ad.Action(func(collection *mgo.Collection) error {

		if err := collection.EnsureIndex(mgo.Index{
			Key:        []string{"email"},
			Unique:     true,
			DropDups:   true,
			Background: true,
			Sparse:     true,
		}); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	return &UserService{app.Logger, ad, app.EmailSender}
}
