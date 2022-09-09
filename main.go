package main

import (
	"knocker/1009/handler"
)

type Note struct {
	id     int
	name   string
	text   string
	access []int
	ttl    int
}

func main() {
	//Notes := make(map[string]Note)
	//Users := make(map[string]string)
	handler.HandleRequest()
}
