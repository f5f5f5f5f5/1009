package service

type Note struct {
	Id     int      `json:"Id"`
	Name   string   `json:"Name"`
	Text   string   `json:"Text"`
	Access []string `json:"Access"`
	Ttl    string   `json:"Ttl"`
}

var NoteRange = []Note{}
