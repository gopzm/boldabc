package students

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"pzm/boldabc/environment"

	"github.com/gorilla/mux"
)

var updateTmpl = template.Must(template.ParseFiles("templates/students/update.html"))

func queryStudent(env *environment.Env, id int64) (*Student, error) {
	db, err := env.OpenDB()
	if err != nil {
		return nil, fmt.Errorf("Cannot open Database: %v", err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM students WHERE id=?", id)
	if err != nil {
		return nil, fmt.Errorf("Database query error: %v", err)
	}
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
				return nil, fmt.Errorf("Cannot parse Database result: %v", err)
			}
		}
		s.Id = id
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

		return &s, nil
	}
	return nil, fmt.Errorf("No student found for student id: %d", id)

}

func UpdateHandler(env *environment.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ok := env.CheckLoggedIn(w, r)
		if !ok {
			return
		}
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid Student Id: %s", vars["id"]), http.StatusBadRequest)
			return
		}
		if r.Method != http.MethodPost {
			s, err := queryStudent(env, id)
			if err != nil {
				http.Error(w, fmt.Sprintf("No student found: %v", err), http.StatusInternalServerError)
				return
			}
			updateTmpl.Execute(w, s)
			return
		}
		s := Student{
			FirstName:    r.FormValue("firstName"),
			LastName:     r.FormValue("lastName"),
			Gender:       r.FormValue("gender"),
			School:       r.FormValue("school"),
			City:         r.FormValue("city"),
			EnglishLevel: r.FormValue("englishLevel"),
		}
		// TODO: parent ID.
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
		stmt, err := db.Prepare("UPDATE students SET firstname=?,lastname=?,gender=?,age=?,grade=?,school=?,city=?,english_level=?,time_created=? WHERE id=?")
		if err != nil {
			http.Error(w, fmt.Sprintf("SQL Error: %v", err), http.StatusInternalServerError)
			return
		}
		_, err = stmt.Exec(s.FirstName, s.LastName, s.Gender, s.Age, s.Grade, s.School, s.City, s.EnglishLevel, s.TimeCreated, id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Database Error: %v", err), http.StatusInternalServerError)
			return
		}
		updateTmpl.Execute(w, s)
	}
}
