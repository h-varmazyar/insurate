package plate

import (
	"context"
	gormext "github.com/h-varmazyar/gopack/gorm"
	personRepo "github.com/h-varmazyar/insurate/internal/core/repository/person"
)

type Plate struct {
	gormext.UniversalModel
	Person      *personRepo.Person
	Alphabet    string
	StartNumber int8
	EndNumber   int8
	RegionCode  int8
	Text        string
}

type Repository interface {
	ReturnByText(ctx context.Context, text string) (*Plate, error)
	Create(ctx context.Context, plate *Plate) error
}
