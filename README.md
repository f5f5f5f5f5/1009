# This is still not done, I`m currently adding some functions

## What is this?
##### This is a simplified note service
![home](https://user-images.githubusercontent.com/107932413/189506223-4ddf38b1-b2af-4e6b-a3d1-091021132eda.png)

## What is realised now?
##### Registration
![reg](https://user-images.githubusercontent.com/107932413/189506237-f22cd565-8784-431e-9519-0320d2a8ade8.png)
Registration is quite simple. It parses values from html forms, checks passwords to match,
generates empty slice of structers (which contains username and password), unmarshals json Users file,
checks login to not match with one that already exists, adds a structure with a new username and password to an empty slice, and then appends this slice to a Users slice and marshals it back to Json.
#####Authorisation
![auth](https://user-images.githubusercontent.com/107932413/189506261-23e33f8a-b432-466c-bb55-2d06ae09f44a.png)
Authorisation gets html forms, checks it is not empty, if authorisation passes it creates a new session and adds cookies.
#####Adding notes
![addnote](https://user-images.githubusercontent.com/107932413/189506511-352f359f-4cb0-4f42-82de-284fb9425931.png)
This function simply adds given name, text and ttl of a note to a json file, also it sets access by default to a current user using cookie.
#####Check authorisation
This function deny access to a chosen pages if user is not authorised to a system
```
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
```

##Work is still in progress on
#####Editing and deleting notes
#####Deleting note once its time is up



