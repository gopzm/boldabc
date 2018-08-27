// Package students handles strudents related logic.
package students

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"pzm/boldabc/environment"
)

var listTmpl = template.Must(template.ParseFiles("templates/students/list.html"))

type ListPage struct {
	Students []*Student
}

func ListHandler(env *environment.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, err := env.OpenDB()
		if err != nil {
			http.Error(w, fmt.Sprintf("Cannot open Database: %v", err), http.StatusInternalServerError)
			return
		}
		defer db.Close()
		rows, err := db.Query("SELECT * FROM students")
		if err != nil {
			http.Error(w, fmt.Sprintf("Database query error: %v", err), http.StatusNotFound)
			return
		}

		var p ListPage
		for rows.Next() {
			var s Student
			var (
				pid    sql.NullInt64
				fn     sql.NullString
				ln     sql.NullString
				gender sql.NullString
				age    sql.NullInt64
				grade  sql.NullInt64
				school sql.NullString
				city   sql.NullString
				el     sql.NullString
				tc     sql.NullInt64
			)
			if err := rows.Scan(&s.Id, &pid, &fn, &ln, &gender,
				&age, &grade, &school, &city, &el, &tc); err != nil {
				if err != nil {
					http.Error(w, fmt.Sprintf("Cannot parse Database result: %v", err), http.StatusInternalServerError)
					return
				}
			}
			s.ParentId = pid.Int64
			s.FirstName = fn.String
			s.LastName = ln.String
			s.Gender = gender.String
			s.Age = age.Int64
			s.Grade = grade.Int64
			s.School = school.String
			s.City = city.String
			s.EnglishLevel = el.String
			s.TimeCreated = tc.Int64

			p.Students = append(p.Students, &s)
		}
		ok := env.CheckLoggedIn(w, r)
		if !ok {
			return
		}
		listTmpl.Execute(w, p)
	}
}
