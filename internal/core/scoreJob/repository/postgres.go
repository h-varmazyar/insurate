package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
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
		return nil, errors.New("invalid db instance in score job")
	}
	return &postgresRepository{
		db:     db,
		logger: logger,
	}, nil
}

func (r *postgresRepository) Create(_ context.Context, job *entity.ScoreJob) error {
	if err := r.db.Create(job).Error; err != nil {
		return err
	}
	return nil
}

func (r *postgresRepository) Status(ctx context.Context, jobId uuid.UUID) (entity.JobStatus, error) {
	return 0, nil
}
