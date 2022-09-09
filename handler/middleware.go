package handler

import (
	"log"
	"net/http"
	"time"
)

func checkAuth(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(manager.cookieName)
		if err != nil {
			log.Printf("Can not find cookie while auth check: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		sid := cookie.Value
		session, err := manager.SessionRead(sid)
		if err != nil {
			log.Printf("%v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if session.Created.Add(time.Second * time.Duration(manager.maxlifetime)).Before(time.Now()) {
			manager.sessionDestroy(sid)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		log.Printf("Middleware auth passed")
		handler(w, r)
	}
}
