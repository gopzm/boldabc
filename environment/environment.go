// environment is the server Environment.
package environment

import (
	"database/sql"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
)

type Env struct {
	store          *sessions.CookieStore
	driverName     string
	dataSourceName string
}

func New(store *sessions.CookieStore, driverName, dataSourceName string) *Env {
	return &Env{
		store:          store,
		driverName:     driverName,
		dataSourceName: dataSourceName,
	}
}

func (e *Env) isLoggedIn(r *http.Request) bool {
	session, err := e.store.Get(r, "cookie-name")
	if err != nil {
		return false
	}
	// Check if user is authenticated
	auth, ok := session.Values["authenticated"].(bool)
	return ok && auth
}

func (e *Env) CheckLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	if !e.isLoggedIn(r) {
		http.Redirect(w, r, "/", 302)
		return false
	}
	return true
}

func (e *Env) OpenDB() (*sql.DB, error) {
	return sql.Open(e.driverName, e.dataSourceName)
}
