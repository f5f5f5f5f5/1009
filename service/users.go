package service

type UserUP struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

type Settings struct {
	Users []UserUP
}

const UsersFilename = "1009/service/Users.json"
