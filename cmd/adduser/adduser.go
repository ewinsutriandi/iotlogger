package main

import (
	"errors"
	"os"

	"github.com/ewinsutriandi/iotlogger"
)

func main() {
	ctl, err := iotlogger.NewController()
	if err == nil {
		if len(os.Args) == 3 {
			username := os.Args[1]
			password := os.Args[2]
			ctl.AddUser(username, password)
		} else {
			err = errors.New("Usage :adduser [username] [password]")
		}
	}
	iotlogger.HandleError(err)
}
