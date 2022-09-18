package handler

import (
	//"crypto/tls"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"knocker/1009/service"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var manager = NewManager("session_id", 60*10)

func home_page(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/home.html", "templates/top.html", "templates/bot.html")
	if err != nil {
		log.Fatal(err)
	}

	var allNotes []service.Note

	const UsersFilename = "service/Notes.json"

	rawDataIn, err := ioutil.ReadFile(UsersFilename)
	if err != nil {
		log.Printf("Cannot load file: %v", err)
	}

	err = json.Unmarshal(rawDataIn, &allNotes)
	if err != nil {
		log.Printf("Failed to unmarshall with error: %v", err)
	}

	var homeNotes []service.Note

	currentCookie, err := r.Cookie(manager.cookieName)
	if err != nil {
		log.Printf("Failed to get current cookie from cookie name: %v", err)
	}

	sid := currentCookie.Value
	currentSession, err := manager.SessionRead(sid)
	if err != nil {
		log.Printf("Failed to get user id from session: %v", err)
	}

	userid := currentSession.Login

	for _, value := range allNotes {
		for _, valeu := range value.Access {
			if valeu == userid {
				homeNotes = append(homeNotes, value)
			}
		}
	}

	tmpl.ExecuteTemplate(w, "home", homeNotes)
}

func login_page(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/login.html", "templates/top.html", "templates/bot.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.ExecuteTemplate(w, "login", nil)
}

func checkin(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		http.ServeFile(w, r, "templates/login.html")
		return
	}

	login := r.FormValue("login")
	password := r.FormValue("password")

	login, password = service.CheckUP(login, password)

	if (login == "") || (password == "") {
		log.Printf("Auth failed")
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else {
		log.Printf("Auth passed")
		sid := manager.sessionId()
		session, err := manager.SessionInit(sid)
		if err != nil {
			log.Printf("Can not create session: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		session.Login = login
		session.Created = time.Now()
		cookie := http.Cookie{
			Name:  manager.cookieName,
			Value: sid,
			Path:  "/",
		}
		log.Printf("Create cookie: %+v", cookie)
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/home/", http.StatusSeeOther)
	}
}

func registration_page(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/registration.html", "templates/top.html", "templates/bot.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.ExecuteTemplate(w, "registration", nil)
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "templates/registration.html")
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("Can not parse form Sign Up: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	login := r.FormValue("login")
	password := r.FormValue("password")
	password1 := r.FormValue("password1")

	if password != password1 {
		log.Printf("Passwords do not match")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sv := service.UserUP{
		Username: login,
		Password: password,
	}
	var vs []service.UserUP

	const UsersFilename = "service/Users.json"

	rawDataIn, err := ioutil.ReadFile(UsersFilename)
	if err != nil {
		log.Printf("Cannot load file: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(rawDataIn, &vs)
	if err != nil {
		log.Printf("Failed to unmarshall with error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, value := range vs {
		if (value.Username) == login {
			log.Println("This name was already taken")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	vs = append(vs, sv)

	boolVar, err := json.Marshal(vs)

	if err != nil {
		log.Printf("Json marshalling failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = ioutil.WriteFile(UsersFilename, boolVar, 0)

	if err != nil {
		log.Printf("Cannot write updated Users file: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/home/", http.StatusSeeOther)
}

func newnote_page(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/newnote.html", "templates/top.html", "templates/bot.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.ExecuteTemplate(w, "newnote", nil)
}

func save_note(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "templates/login.html")
		return
	}

	name := r.FormValue("Name")
	text := r.FormValue("Text")
	ttl := r.FormValue("Ttl")

	currentCookie, err := r.Cookie(manager.cookieName)
	if err != nil {
		log.Printf("Failed to get current cookie from cookie name: %v", err)
	}
	sid := currentCookie.Value
	currentSession, err := manager.SessionRead(sid)
	if err != nil {
		log.Printf("Failed to get user id from session: %v", err)
	}

	userid := currentSession.Login

	service.NewNote(userid, name, text, ttl)
	http.Redirect(w, r, "/home/", http.StatusSeeOther)
}

func editNote(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	strId := r.FormValue("Id")
	name := r.FormValue("Name")
	text := r.FormValue("Text")
	strTtl := r.FormValue("Ttl")

	id, err := strconv.Atoi(strId)
	if err != nil {
		log.Printf("Failed to convert id to int: %v", err)
	}

	ttl, err := strconv.Atoi(strTtl)
	if err != nil {
		log.Printf("Failed to convert ttl to int: %v", err)
	}

	currentCookie, err := r.Cookie(manager.cookieName)
	if err != nil {
		log.Printf("Failed to get current cookie from cookie name: %v", err)
	}

	sid := currentCookie.Value

	currentSession, err := manager.SessionRead(sid)
	if err != nil {
		log.Printf("Failed to get user id from session: %v", err)
	}

	userid := currentSession.Login

	service.EditNote(name, text, userid, id, ttl)

	http.Redirect(w, r, "/yournotes/", http.StatusSeeOther)
}

func userNotes_page(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/usernotes.html", "templates/top.html", "templates/bot.html")
	if err != nil {
		log.Fatal(err)
	}

	var allNotes []service.Note

	const notesFilename = "service/Notes.json"

	rawDataIn, err := ioutil.ReadFile(notesFilename)
	if err != nil {
		log.Printf("Cannot load file: %v", err)
	}

	err = json.Unmarshal(rawDataIn, &allNotes)
	if err != nil {
		log.Printf("Failed to unmarshall with error: %v", err)
	}

	var userNotes []service.Note

	currentCookie, err := r.Cookie(manager.cookieName)
	if err != nil {
		log.Printf("Failed to get current cookie from cookie name: %v", err)
	}

	sid := currentCookie.Value

	currentSession, err := manager.SessionRead(sid)
	if err != nil {
		log.Printf("Failed to get user id from session: %v", err)
	}

	userid := currentSession.Login

	for _, value := range allNotes {
		for _, valeu := range value.Access {
			if valeu == userid {
				userNotes = append(userNotes, value)
			}
		}
	}

	tmpl.ExecuteTemplate(w, "usernotes", userNotes)
}

func editNote_page(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strid := vars["Id"]

	id, err := strconv.Atoi(strid)

	if err != nil {
		log.Print("Failed to convert id to int: editNote_page")
	}

	var allNotes []service.Note

	const notesFilename = "service/Notes.json"

	rawDataIn, err := ioutil.ReadFile(notesFilename)
	if err != nil {
		log.Printf("Cannot load file: %v", err)
	}

	err = json.Unmarshal(rawDataIn, &allNotes)
	if err != nil {
		log.Printf("Failed to unmarshall with error: %v", err)
	}

	noteToEdit := service.Note{}

	for _, value := range allNotes {
		if id == value.Id {
			noteToEdit = value
		}
	}
	log.Print(noteToEdit)

	tmpl, err := template.ParseFiles("templates/editnote.html", "templates/top.html", "templates/bot.html")
	if err != nil {
		log.Fatal(err)
	}

	tmpl.ExecuteTemplate(w, "editnote", noteToEdit)

}

func deleteNote(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	strid := vars["Id"]

	id, err := strconv.Atoi(strid)
	if err != nil {
		log.Printf("Failed to convert id to int: %v", err)
	}

	currentCookie, err := r.Cookie(manager.cookieName)
	if err != nil {
		log.Printf("Failed to get current cookie from cookie name: %v", err)
	}

	sid := currentCookie.Value

	currentSession, err := manager.SessionRead(sid)
	if err != nil {
		log.Printf("Failed to get user id from session: %v", err)
	}

	userid := currentSession.Login

	service.DeleteNote(id, userid)

	http.Redirect(w, r, "/yournotes/", http.StatusSeeOther)
}

func HandleRequest() {
	router := mux.NewRouter()
	router.HandleFunc("/", home_page)
	router.HandleFunc("/home/", home_page)
	router.HandleFunc("/newnote/", checkAuth(newnote_page))
	router.HandleFunc("/save_note/", save_note)
	router.HandleFunc("/login/", login_page)
	router.HandleFunc("/checkin/", checkin)
	router.HandleFunc("/registration/", registration_page)
	router.HandleFunc("/register/", register)
	router.HandleFunc("/yournotes/", userNotes_page)
	router.HandleFunc("/edit/{id:[0-9]+}/", editNote).Methods("POST")
	router.HandleFunc("/delete/{id:[0-9]+}/", deleteNote)
	router.HandleFunc("/edit/{id:[0-9]+}/", editNote_page).Methods("GET")
	http.Handle("/", router)
	http.ListenAndServe(":5040", nil)
}
