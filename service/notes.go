package service

type Note struct {
	UserId int    `json:"UserId"`
	Id     int    `json:"Id"`
	Name   string `json:"Name"`
	Text   string `json:"Text"`
	Access []int  `json:"Access"`
	Ttl    string `json:"Ttl"`
}

var NoteRange = []Note{}
