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
	login := r.FormValue("login")
	password := r.FormValue("password")

	login, password = service.CheckUP(login, password)

	if (login != "") && (password != "") { //switch case empty login/password
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
	tmpl, err := template.ParseFiles("templates/registration.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(w, nil)
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
		log.Printf("Failer to unmarshall with error: %v", err)
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
	tmpl, err := template.ParseFiles("templates/newnote.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(w, nil)
}

func save_note(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("Name")
	text := r.FormValue("Text")
	ttl := r.FormValue("Ttl")
	userid := "Beta"
	// if err != nil {
	// 	log.Printf("Failed to get user id from session: %v", err)
	// }
	service.NewNote(userid, name, text, ttl)
	http.Redirect(w, r, "/home/", http.StatusSeeOther)
}

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
	// //srv := &http.Server{
	// 	Addr:    "localhost:5040",
	// 	Handler: mux,
	// 	// TLSConfig: &tls.Config{
	// 	// 	MinVersion:               tls.VersionTLS13,
	// 	// 	PreferServerCipherSuites: true,
	// 	// },
	// }
	http.Handle("/", mux)
	http.ListenAndServe(":5040", nil)
	// err := srv.ListenAndServeTLS("1009/key/server.crt", "1009/key/server.key")
	// if err != nil {
	// 	log.Fatalf("can not listen port 5040: %v", err)
	// }
}
