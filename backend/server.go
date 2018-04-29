package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	store  = sessions.NewCookieStore(securecookie.GenerateRandomKey(32))
	router = mux.NewRouter()

	validUsers = map[string]string{
		"pzm": "IAmTheMaster",
		"sy":  "IAmThePig",
	}
)

// index page

func indexPageHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		fmt.Fprintf(w, loginPage)
		return
	}
	userName, ok := session.Values["username"].(string)
	if !ok {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	fmt.Fprintf(w, internalPage, userName)
}

const loginPage = `
<h1>Login</h1>
<form method="post" action="/login">
    <label for="username">User name</label>
    <input type="text" id="username" name="username">
    <label for="password">Password</label>
    <input type="password" id="password" name="password">
    <button type="submit">Login</button>
</form>
`

// internal page

const internalPage = `
<h1>Internal</h1>
<hr>
<small>User: %s</small>
<form method="post" action="/logout">
    <button type="submit">Logout</button>
</form>
`

// login handler

func loginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	passwd := r.FormValue("password")
	if username != "" && passwd != "" && validUsers[username] == passwd {
		session, _ := store.Get(r, "cookie-name")
		session.Values["username"] = username
		session.Values["authenticated"] = true
		session.Save(r, w)
	}
	http.Redirect(w, r, "/", 302)
}

// logout handler

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	session.Values["authenticated"] = false
	session.Save(r, w)
	http.Redirect(w, r, "/", 302)
}

func main() {
	router.HandleFunc("/", indexPageHandler)
	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/logout", logoutHandler).Methods("POST")

	http.ListenAndServe(":8080", router)
}
