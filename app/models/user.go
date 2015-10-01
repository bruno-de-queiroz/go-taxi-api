package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Facebook struct {
	ID                string `bson:"id"`
	NickName          string
	Description       string
	AvatarURL         string
	Location          string
	AccessToken       string
	AccessTokenSecret string
	RawData           map[string]interface{}
}

type User struct {
	ID          bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Email       string        `json:"email"`
	Name        string        `json:"name"`
	Password    string        `json:"-"`
	AccessToken *string       `json:"-" bson:"accessToken,omitempty"`
	ResetToken  *string       `json:"-" bson:"resetToken,omitempty"`
	Facebook    *Facebook     `json:"-" bson:"facebook,omitempty"`
	Created     time.Time     `json:"created_at"`
	Updated     time.Time     `json:"updated_at"`
}
