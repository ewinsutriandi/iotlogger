package main

import (
	"errors"
	"log"
	"os"

	"strconv"

	"github.com/ewinsutriandi/iotlogger"
)

func main() {
	var err error
	log.Printf("arguments length = %d", len(os.Args))
	if len(os.Args) >= 2 {
		command := os.Args[1]
		switch command {
		case "adduser":
			err = addUser(os.Args)
		case "checkuser":
			err = checkUser(os.Args)
		case "addtelem":
			err = addTelem(os.Args)
		default:
			err = errors.New("Unknown command")
		}

	} else {
		err = errors.New("Please specify command")
	}
	iotlogger.HandleError(err)
}
func addUser(args []string) (err error) {
	ctl, err := iotlogger.NewController()
	if err == nil {
		if len(os.Args) == 4 {
			username := os.Args[2]
			password := os.Args[3]
			err = ctl.AddUser(username, password)
		} else {
			err = errors.New("Usage :user add [username] [password]")
		}
	}
	return
}
func checkUser(args []string) (err error) {
	ctl, err := iotlogger.NewController()
	if err == nil {
		if len(os.Args) == 3 {
			username := os.Args[2]
			user, err := ctl.GetUser(username)
			if err == nil {
				log.Printf("found user :%s with id %s", user.Username, user.ID)
			} else {
				log.Printf("user not found")
			}
		} else {
			err = errors.New("Usage :user check [username]")
		}
	}
	return
}
func addTelem(args []string) (err error) {
	ctl, err := iotlogger.NewController()
	if err == nil {
		if len(os.Args) == 5 {

			metric := os.Args[2]
			svalue := os.Args[3]
			value, err := strconv.ParseFloat(svalue, 64)
			if err != nil {
				err = errors.New("invalid value, use number as value")
				return err
			}
			unit := os.Args[4]
			telem := iotlogger.NewTelemetryReport(
				"127.0.0.1",
				"local test",
				metric,
				unit, value)
			ctl.AddTelemetryReport(*telem)

		} else {
			err = errors.New("Usage :user addtelem [metric] [value] [unit]")
		}
	}
	return
}
