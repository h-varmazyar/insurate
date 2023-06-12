package score

import (
	"context"
	"github.com/google/uuid"
	gormext "github.com/h-varmazyar/gopack/gorm"
)

type Status int8

const (
	Pending Status = iota
	PreparingData
	Calculating
	Done
	Failed
)

type Score struct {
	gormext.UniversalModel
	NationalCode string
	Value        float64
	Status       Status
}

type Repository interface {
	Create(ctx context.Context, score *Score) error
	Update(ctx context.Context, score *Score) error
	UpdateStatus(ctx context.Context, scoreID uuid.UUID, newStatus Status) error
	Last(ctx context.Context, nationalCode string) (*Score, error)
}
