package handler

import (
	"fmt"
	"html/template"
	"knocker/1009/service"
	"log"
	"net/http"
)

func home_page(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(w, nil)
}

func login_page(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(w, nil)
}

func checkin(w http.ResponseWriter, r *http.Request) {
	// login := r.FormValue("login")
	// password := r.FormValue("password")
	// if login == "admin" && password == "admin" {
	// 	http.Redirect(w, r, "/home/", http.StatusSeeOther)
	// }
}

func registration_page(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/registration.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(w, nil)
}

func register(w http.ResponseWriter, r *http.Request) {
	// regisration actions
}

func newnote_page(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/newnote.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(w, nil)
}

func save_note(w http.ResponseWriter, r *http.Request) {
	var NewNote service.Note
	NewNote.Id = 1
	NewNote.Name = "Pervaya"
	NewNote.Text = "Text"
	NewNote.Access = append(NewNote.Access, 123)
	NewNote.Ttl = "321"
	service.NoteRange = append(service.NoteRange, NewNote)
	fmt.Println(service.NoteRange)
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
