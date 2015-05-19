package models

import "time"

type Application struct {
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	SSN       string    `json:"ssn"`
	Email     string    `json:"email"`
	Created   time.Time `json:"created_at"`
	Updated   time.Time `json:"updated_at"`
}

// a.Valid() is this a validly filled in application?
func (a *Application) Valid() bool {
	return true
}
