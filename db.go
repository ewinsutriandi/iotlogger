package iotlogger

import (
	"errors"
	"log"

	"time"

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

// TelemetryReports struct
// e.g : ID, x.x.x.x, weather station, wind speed, 19, m/s
type TelemetryReports struct {
	ID         bson.ObjectId `json:"id" bson:"_id,omitempty"`
	DeviceIP   string
	DeviceName string
	Datetime   time.Time
	Metric     string
	Value      float64
	Unit       string
}

func NewTelemetryReport(ip, devname, metric, unit string, value float64) *TelemetryReports {
	telem := new(TelemetryReports)
	telem.DeviceIP = ip
	telem.DeviceName = devname
	telem.Metric = metric
	telem.Value = value
	telem.Unit = unit
	telem.Datetime = time.Now()
	return telem
}

// GetUser : get user with specified username from DB
func (c *Controller) GetUser(username string) (user User, err error) {
	defer c.session.Close()
	coll := c.session.DB("iot").C("users")
	err = coll.Find(bson.M{"username": username}).One(&user)
	HandleError(err)
	return user, err
}

// AddUser : add new user to DB
func (c *Controller) AddUser(username, password string) (err error) {
	defer c.session.Close()
	coll := c.session.DB("iot").C("users")
	salt, err := generateRandomSalt(17)
	pass, err := encrypt(password, salt)
	i := bson.NewObjectId()
	cnt, err := coll.Find(bson.M{"username": username}).Count()
	if cnt == 0 {
		err = coll.Insert(&User{i, username, pass, salt})
		if err == nil {
			log.Printf("user :" + username + ", inserted to db ")
		}
	} else {
		err = errors.New("username already exists")
	}
	HandleError(err)
	return
}

//AddTelemetryReport : save TelemetryReports to DB
func (c *Controller) AddTelemetryReport(tel TelemetryReports) (err error) {
	defer c.session.Close()
	coll := c.session.DB("iot").C("telemreports")
	i := bson.NewObjectId()
	tel.ID = i
	err = coll.Insert(tel)
	if err == nil {
		log.Printf("New TelemetryReports inserted to db ")
	}
	HandleError(err)
	return
}
