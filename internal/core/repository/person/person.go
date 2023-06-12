package person

import (
	"context"
	"time"
)

type Gender int8

const (
	Unknown Gender = iota
	Men
	Women
)

type Person struct {
	NationalCode string
	BirthDate    time.Time
	FirstName    string
	LastName     string
	Mobile       string
	Gender
}

type Repository interface {
	Create(ctx context.Context, person *Person) error
	Return(ctx context.Context, nationalCode string) (*Person, error)
}
