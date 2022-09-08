package handler

import (
	"html/template"
	"log"
	"net/http"
)

func home_page(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("1009/templates/home.html", "1009/templates/head.html", "1009/templates/footer.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.ExecuteTemplate(w, "home", nil)
}

func login_page(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("1009/templates/login.html", "1009/templates/head.html", "1009/templates/footer.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.ExecuteTemplate(w, "login", nil)
}

func checkin(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	password := r.FormValue("password")
	if login == "admin" && password == "admin" {
		http.Redirect(w, r, "/home/", http.StatusSeeOther)
	}
}

func home_TEST(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("1009/templates/hz.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(w, nil)
}

func HandleRequest() {
	http.HandleFunc("/", home_page)
	http.HandleFunc("/home/", home_TEST)
	http.HandleFunc("/login/", login_page)
	http.HandleFunc("/checkin/", checkin)
	http.ListenAndServe(":5040", nil)
}
