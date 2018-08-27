package students

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"pzm/boldabc/environment"
)

var addTmpl = template.Must(template.ParseFiles("templates/students/add.html"))

func AddHandler(env *environment.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ok := env.CheckLoggedIn(w, r)
		if !ok {
			return
		}
		if r.Method != http.MethodPost {
			addTmpl.Execute(w, nil)
			return
		}
		s := Student{
			FirstName:    r.FormValue("firstName"),
			LastName:     r.FormValue("lastName"),
			Gender:       r.FormValue("gender"),
			School:       r.FormValue("school"),
			City:         r.FormValue("city"),
			EnglishLevel: r.FormValue("englishLevel"),
			TimeCreated:  time.Now().Unix(),
		}
		age, err := strconv.ParseInt(r.FormValue("age"), 10, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid Age Value: %v", err), http.StatusBadRequest)
			return
		}
		s.Age = age
		grade, err := strconv.ParseInt(r.FormValue("grade"), 10, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid Grade Value: %v", err), http.StatusBadRequest)
			return
		}
		s.Grade = grade
		db, err := env.OpenDB()
		if err != nil {
			http.Error(w, fmt.Sprintf("Cannot open Database: %v", err), http.StatusInternalServerError)
			return
		}
		defer db.Close()
		stmt, err := db.Prepare("INSERT students SET firstname=?,lastname=?,gender=?,age=?,grade=?,school=?,city=?,english_level=?,time_created=?")
		if err != nil {
			http.Error(w, fmt.Sprintf("SQL Error: %v", err), http.StatusInternalServerError)
			return
		}
		_, err = stmt.Exec(s.FirstName, s.LastName, s.Gender, s.Age, s.Grade, s.School, s.City, s.EnglishLevel, s.TimeCreated)
		if err != nil {
			http.Error(w, fmt.Sprintf("Database Error: %v", err), http.StatusInternalServerError)
			return
		}
		addTmpl.Execute(w, struct{ Success bool }{true})
	}
}
