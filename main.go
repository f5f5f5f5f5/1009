package main

import (
	"awesomeProject3/1009/handler"
)

type Note struct {
	id     int
	name   string
	text   string
	access []int
	ttl    int
}

func main() {
	Notes := make(map[string]Note)
	//Users := make(map[string]string)
	var a Note
	a.name = "Nazvanie zametOCHKi"
	a.text = "Ya ebal amerov"
	a.access = append(a.access, 212)
	a.ttl = 112
	a.id = 212
	Notes["12131"] = a
	handler.HandleRequest()
}
