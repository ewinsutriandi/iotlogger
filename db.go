package iotlogger

import (
	"errors"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	MongoServerAddr = "127.0.0.1"
)

type Controller struct {
	session *mgo.Session
}

func NewController() (*Controller, error) {
	session, err := mgo.Dial(MongoServerAddr)
	if err != nil {
		HandleError(err)
		return nil, err
	}
	return &Controller{
		session: session,
	}, nil
}

// User struct
type User struct {
	ID       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Username string
	Password []byte
	Salt     []byte
}

// GetUser : get user with specified username from DB
func (c *Controller) GetUser(username string) (User, error) {
	defer c.session.Close()
	var user User
	coll := c.session.DB("iot").C("users")
	err := coll.Find(bson.M{"username": username}).One(&user)
	if err == nil {
		return user, nil
	}
	HandleError(err)
	return user, err
}

// AddUser : add new user to DB
func (c *Controller) AddUser(username, password string) error {
	defer c.session.Close()
	var err error
	coll := c.session.DB("iot").C("users")
	salt, err := generateRandomSalt(17)
	pass, err := encrypt(password, salt)
	i := bson.NewObjectId()
	cnt, err := coll.Find(bson.M{"username": username}).Count()
	if cnt == 0 {
		err = coll.Insert(&User{i, username, pass, salt})
	} else {
		err = errors.New("username already exists")
	}
	if err != nil {
		HandleError(err)
	}
	return err
}
