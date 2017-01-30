package iotlogger

import (
	"bytes"
	"crypto/rand"

	"golang.org/x/crypto/scrypt"
)

func generateRandomSalt(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func encrypt(pass string, salt []byte) ([]byte, error) {
	dk, err := scrypt.Key([]byte(pass), salt, 16384, 8, 1, 32)
	return dk, err
}

//Authenticate user
func Authenticate(username, pass string) bool {
	//savedPass = get []byte from db
	auth := false
	savedPass, salt := getPassSalt(username)
	key, err := encrypt(pass, salt)
	if err == nil {
		i := bytes.Compare(savedPass, key)
		if i == 0 {
			auth = true
		}
	}
	return auth
}

func getPassSalt(username string) ([]byte, []byte) {
	//get pass & salt from db where username = user
	var err error
	ctl, err := NewController()
	if err == nil {
		User, err := ctl.GetUser(username)
		if err == nil {
			return []byte(User.Password), []byte(User.Salt)
		}
	}
	panic(err)
}
