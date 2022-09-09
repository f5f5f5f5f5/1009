package handler

import (
	"awesomeProject3/1009/service"
	"crypto/tls"
	"html/template"
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
	} else {
		log.Printf("Auth failed")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
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
	mux := http.NewServeMux()
	mux.HandleFunc("/", home_page)
	mux.HandleFunc("/home/", checkAuth(home_page))
	mux.HandleFunc("/newnote/", newnote_page)
	mux.HandleFunc("/save_note/", save_note)
	mux.HandleFunc("/login/", login_page)
	mux.HandleFunc("/checkin/", checkin) // действия с авторизацией
	mux.HandleFunc("/registration/", registration_page)
	mux.HandleFunc("/register/", register) // действия с регистрацией
	srv := &http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS13,
			PreferServerCipherSuites: true,
		},
	}
	err := srv.ListenAndServeTLS("1009/key/server.crt", "1009/key/server.key")
	if err != nil {
		log.Fatalln("can not listen port 8080: %v", err)
	}
}
