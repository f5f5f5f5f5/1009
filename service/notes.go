package service

import (
	"encoding/json"
	"io/ioutil"
	"log"

	//"time"
	"strconv"
)

type Note struct {
	UserId string   `json:"UserID"`
	Id     int      `json:"Id"`
	Name   string   `json:"Name"`
	Text   string   `json:"Text"`
	Access []string `json:"Access"`
	Ttl    int      `json:"Ttl"`
}

const notesFilename = "service/Notes.json"

func NewNote(userid, name, text, ttl string) {
	var noteArray []Note

	rawDataIn, err := ioutil.ReadFile(notesFilename)
	if err != nil {
		log.Printf("Cannot load file: %v", err)
	}

	err = json.Unmarshal(rawDataIn, &noteArray)
	if err != nil {
		log.Printf("Failed to unmarshall with error: %v", err)
	}

	var accessArray []string
	accessArray = append(accessArray, userid)

	intttl, err := strconv.Atoi(ttl)

	if err != nil {
		log.Printf("Failed to convert time: %v", err)
	}

	idNote := noteArray[(len(noteArray) - 1)]

	newnote := Note{
		UserId: userid,
		Id:     idNote.Id + 1,
		Name:   name,
		Text:   text,
		Access: accessArray,
		Ttl:    intttl,
	}

	noteArray = append(noteArray, newnote)

	boolVar, err := json.Marshal(noteArray)

	if err != nil {
		log.Printf("Json marshalling failed: %v", err)
	}

	err = ioutil.WriteFile(notesFilename, boolVar, 0)

	if err != nil {
		log.Printf("Cannot write updated Notes file: %v", err)
	}
}

func DeleteNote(id int, userid string) {
	var allNotes []Note

	rawDataIn, err := ioutil.ReadFile(notesFilename)
	if err != nil {
		log.Printf("Cannot load file: %v", err)
	}

	err = json.Unmarshal(rawDataIn, &allNotes)
	if err != nil {
		log.Printf("Failed to unmarshall with error: %v", err)
	}

	for i, value := range allNotes {
		if value.Id == id {
			if value.UserId == userid {
				allNotes = removeByIndex(allNotes, i)
			}
		}
	}

	boolVar, err := json.Marshal(allNotes)

	if err != nil {
		log.Printf("Json marshalling failed: %v", err)
	}

	err = ioutil.WriteFile(notesFilename, boolVar, 0)

	if err != nil {
		log.Printf("Cannot write updated Notes file: %v", err)
	}

}

func EditNote(name, text, userid, access string, id, ttl int) {
	var allNotes []Note

	rawDataIn, err := ioutil.ReadFile(notesFilename)
	if err != nil {
		log.Printf("Cannot load file: %v", err)
	}

	err = json.Unmarshal(rawDataIn, &allNotes)
	if err != nil {
		log.Printf("Failed to unmarshall with error: %v", err)
	}

	for i, value := range allNotes {
		if value.Id == id {
			log.Print("got a match of IDs")
			if value.UserId == userid {
				log.Print("got a match of UserIds")
				value.Name = name
				value.Text = text
				value.Ttl = ttl

				if access != "" {
					value.Access = accessCheck(access, value.Access)
				}

				allNotes[i] = value
				log.Print(allNotes)
			} else {
				log.Printf("Note UserId and current session UserId do not match, current user: %s, note userid: %s", userid, value.UserId)
				return
			}
		} else {
			log.Print("There is no note with given Id")
		}
	}

	boolVar, err := json.Marshal(allNotes)

	if err != nil {
		log.Printf("Json marshalling failed: %v", err)
	}

	err = ioutil.WriteFile(notesFilename, boolVar, 0)

	if err != nil {
		log.Printf("Cannot write updated Notes file: %v", err)
	}
}

func removeByIndex(array []Note, index int) []Note {
	return append(array[:index], array[index+1:]...)
}

func accessCheck(access string, accesses []string) []string {
	for _, value := range accesses {
		if value == access {
			return removeAccess(access, accesses)
		}
	}
	return shareAccess(access, accesses)
}

func shareAccess(access string, accesses []string) []string {
	accesses = append(accesses, access)
	return accesses
}

func removeAccess(access string, accesses []string) []string {
	for i, value := range accesses {
		if access == value {
			accesses = append(accesses[:i], accesses[i+1:]...)
		}
	}
	return accesses
}
