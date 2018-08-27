// Package students handles strudents related logic.
package students

// Student is the db model of student for displaying.
type Student struct {
	Id           int64
	ParentId     int64
	FirstName    string
	LastName     string
	Gender       string
	Age          int64
	Grade        int64
	School       string
	City         string
	EnglishLevel string
	TimeCreated  int64
}
