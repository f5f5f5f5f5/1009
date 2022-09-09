package service

type User struct {
	Id       int    `json:"Id"`
	Name     string `json:"Name"`
	Username string `json:"Username"`
	Password string `json:"Password"`
}

var UserRange []User
