package handler

import (
	"html/template"
	"awesomeProject3/1009/service"
	"log"
	"net/http"
	"time"
)

var manager = NewManager("session_id", 60*10)

func home_page(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("1009/templates/home.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(w, nil)
}

func login_page(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("1009/templates/login.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(w, nil)
}

func checkin(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	password := r.FormValue("password")
	if login == "admin" && password == "admin" {
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
			HttpOnly: true,
			Name:     manager.cookieName,
			Value:    sid,
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/home/", http.StatusSeeOther)
}

func registration_page(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("1009/templates/registration.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(w, nil)
}

func register(w http.ResponseWriter, r *http.Request) {
	// regisration actions
}

func newnote_page(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("1009/templates/newnote.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(w, nil)
}

func save_note(w http.ResponseWriter, r *http.Request) {
	var NewNote service.Note
	// NewNote.Id = 1, 2, 3
	NewNote.Name = r.FormValue("Name")
	NewNote.Text = r.FormValue("Text")
	// NewNote.Access = user id
	NewNote.Ttl = r.FormValue("Ttl")
}

func HandleRequest() {
	http.HandleFunc("/", home_page)
	http.HandleFunc("/home/", home_page)
	http.HandleFunc("/newnote/", newnote_page)
	http.HandleFunc("/save_note/", save_note)
	http.HandleFunc("/login/", login_page)
	http.HandleFunc("/checkin/", checkin) // действия с авторизацией
	http.HandleFunc("/registration/", registration_page)
	http.HandleFunc("/register/", register) // действия с регистрацией
	http.ListenAndServe(":5040", nil)
}
