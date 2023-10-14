package repository

import (
	"context"
	"errors"
	"github.com/h-varmazyar/insurate/internal/entity"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type postgresRepository struct {
	db     *gorm.DB
	logger *log.Logger
}

func newPostgresRepository(_ context.Context, logger *log.Logger, db *gorm.DB) (Repository, error) {
	if db == nil {
		return nil, errors.New("invalid db instance in score report")
	}
	return &postgresRepository{
		db:     db,
		logger: logger,
	}, nil
}

func (r *postgresRepository) Create(_ context.Context, job *entity.ScoreReport) error {
	if err := r.db.Create(job).Error; err != nil {
		return err
	}
	return nil
}

func (r *postgresRepository) ReturnByTrackingId(_ context.Context, trackingId string) (*entity.ScoreReport, error) {
	scoreReport := new(entity.ScoreReport)
	err := r.db.Model(new(entity.ScoreReport)).Where("tracking_id = ?", trackingId).First(scoreReport).Error
	if err != nil {
		return nil, err
	}
	return scoreReport, nil
}
