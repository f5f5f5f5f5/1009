package handler

import (
	//"crypto/tls"
	"encoding/json"
	//"fmt"
	"html/template"
	"io/ioutil"
	"knocker/1009/service"
	"log"
	"net/http"
	"time"
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
		log.Printf("Failer to unmarshall with error: %v", err)
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
		http.Redirect(w, r, "/home", http.StatusSeeOther)
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
	}

	err = json.Unmarshal(rawDataIn, &vs)
	if err != nil {
		log.Printf("Failed to unmarshall with error: %v", err)
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
	}

	err = ioutil.WriteFile(UsersFilename, boolVar, 0)

	if err != nil {
		log.Printf("Cannot write updated Users file: %v", err)
	}
}

func newnote_page(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/newnote.html", "templates/top.html", "templates/bot.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.ExecuteTemplate(w, "newnote", nil)
}

func save_note(w http.ResponseWriter, r *http.Request) {
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

func editNote_page(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/editnote.html", "templates/top.html", "templates/bot.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.ExecuteTemplate(w, "editnote", nil)
}

// func deleteNote(w http.ResponseWriter, r *http.Request) {
//
// }

func HandleRequest() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home_page)
	mux.HandleFunc("/home/", home_page)
	mux.HandleFunc("/newnote/", checkAuth(newnote_page))
	mux.HandleFunc("/save_note/", save_note)
	mux.HandleFunc("/login/", login_page)
	mux.HandleFunc("/checkin/", checkin) // действия с авторизацией
	mux.HandleFunc("/registration/", registration_page)
	mux.HandleFunc("/register/", register) // действия с регистрацией
	mux.HandleFunc("/editnote", checkAuth(editNote_page))
	//mux.HandleFunc("/delete", deleteNote)
	http.Handle("/", mux)
	http.ListenAndServe(":5040", nil)
}
