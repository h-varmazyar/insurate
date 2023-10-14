package entity

import (
	gormext "github.com/h-varmazyar/gopack/gorm"
)

type ScoreReport struct {
	gormext.UniversalModel
	TrackingId string
}
