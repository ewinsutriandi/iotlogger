package main

import (
	"errors"
	"log"
	"os"

	"github.com/ewinsutriandi/iotlogger"
)

func main() {
	var err error
	log.Printf("arguments length = %d", len(os.Args))
	if len(os.Args) >= 2 {
		command := os.Args[1]
		switch command {
		case "add":
			err = addUser(os.Args)
		case "check":
			err = checkUser(os.Args)
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
			username := os.Args[1]
			password := os.Args[2]
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
