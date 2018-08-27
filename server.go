package main

import (
	"fmt"
	"net/http"

	"pzm/boldabc/environment"
	"pzm/boldabc/students"

	_ "github.com/go-sql-driver/mysql"
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
<h1>管理员登录</h1>
<form method="post" action="/login">
    <label for="username">管理员</label>
    <input type="text" id="username" name="username">
    <label for="password">密码</label>
    <input type="password" id="password" name="password">
    <button type="submit">登录</button>
</form>
`

// internal page

const internalPage = `
<h1>BoldABC 内部管理页面</h1>
<hr>
<form method="post" action="/logout">
    <label>管理员: %s</label>
    <button type="submit">退出</button>
</form>

<div>
	<h2>功能</h2>
	<ul>
		<li><a href="./students">查看学生列表</a></li>
		<li><a href="./add_student">添加学生</a></li>
	</ul>
</div>
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
	env := environment.New(store, "mysql", "root:@tcp(127.0.0.1:3306)/boldabc")

	router.HandleFunc("/", indexPageHandler)
	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/logout", logoutHandler).Methods("POST")
	router.HandleFunc("/students", students.ListHandler(env))
	router.HandleFunc("/add_student", students.AddHandler(env))
	router.HandleFunc("/update_student/{id}", students.UpdateHandler(env))

	http.ListenAndServe(":8080", router)
}
