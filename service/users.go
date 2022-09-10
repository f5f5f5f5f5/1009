package service

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type UserUP struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

// type Settings struct {
// 	Users []UserUP
// }

func CheckUP(login, password string) (string, string) {
	var vs []UserUP

	const UsersFilename = "service/Users.json"

	rawDataIn, err := ioutil.ReadFile(UsersFilename)
	if err != nil {
		log.Printf("Cannot load Users file: %v", err)
	}

	err = json.Unmarshal(rawDataIn, &vs)
	if err != nil {
		log.Printf("Invalid json format: %v", err)
	}

	for _, value := range vs {
		if (value.Username) == login && (value.Password) == password {
			return login, password
		}
	}
	return "", ""
}
