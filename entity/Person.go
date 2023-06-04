package entity

import "time"

type Person struct {
	NationalCode string
	BirthDate    time.Time
	FirstName    string
	LastName     string
	Mobile       string
}
