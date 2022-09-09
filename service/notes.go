package service

type Note struct {
	Id     string   `json:"Id"`
	Name   string   `json:"Name"`
	Text   string   `json:"Text"`
	Access []string `json:"Access"`
	Ttl    string   `json:"Ttl"`
}

type NoteRange struct {
	UserNotes []Note
}
