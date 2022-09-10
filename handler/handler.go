package handler

import (
	//"crypto/tls"
	"encoding/json"
	"fmt"
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

	if (login != "") && (password != "") {
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
		log.Fatal("Cannot load settings:", err)
	}

	err = json.Unmarshal(rawDataIn, &vs)
	if err != nil {
		log.Fatal("Invalid settings format:", err)
	}

	vs = append(vs, sv)

	boolVar, err := json.Marshal(vs)

	if err != nil {
		log.Fatal("Json marshalling failed:", err)
	}

	err = ioutil.WriteFile(UsersFilename, boolVar, 0)

	if err != nil {
		log.Fatal("Cannot write updated settings file:", err)
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
	var NewNote service.Note
	//NewNote.Id = 1
	NewNote.Name = r.FormValue("Name")
	NewNote.Text = r.FormValue("Text")
	//NewNote.Access = append(NewNote.Access, "login")
	//NewNote.Ttl = r.FormValue("Ttl")
	service.NoteRange = append(service.NoteRange, NewNote)
	fmt.Println(service.NoteRange)
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
