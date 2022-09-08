package handler

import (
	"html/template"
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
	// regisretion actions
}

func HandleRequest() {
	http.HandleFunc("/", home_page)
	http.HandleFunc("/home/", home_page)
	http.HandleFunc("/login/", login_page)
	http.HandleFunc("/checkin/", checkin) // действия с авторизацией
	http.HandleFunc("/registration/", registration_page)
	http.HandleFunc("/register/", register) // действия с регистрацией
	http.ListenAndServe(":5040", nil)
}
