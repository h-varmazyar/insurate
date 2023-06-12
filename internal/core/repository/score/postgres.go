package score

import (
	"context"
	"errors"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type postgresRepository struct {
	db     *gorm.DB
	logger *log.Logger
}

func NewPostgresRepository(_ context.Context, logger *log.Logger, db *gorm.DB) (Repository, error) {
	if db == nil {
		return nil, errors.New("invalid db instance for score")
	}
	return &postgresRepository{
		db:     db,
		logger: logger,
	}, nil
}

func (r *postgresRepository) Create(_ context.Context, score *Score) error {
	if err := r.db.Save(score).Error; err != nil {
		return err
	}
	return nil
}

func (r *postgresRepository) Update(_ context.Context, score *Score) error {
	if err := r.db.Model(new(Score)).Updates(score).Error; err != nil {
		return err
	}
	return nil
}

func (r *postgresRepository) UpdateStatus(_ context.Context, scoreID uuid.UUID, newStatus Status) error {
	if err := r.db.Model(new(Score)).Where("id = ?", scoreID).UpdateColumn("status", newStatus).Error; err != nil {
		return err
	}
	return nil
}

func (r *postgresRepository) Last(_ context.Context, nationalCode string) (*Score, error) {
	score := new(Score)
	if err := r.db.Model(new(Score)).Where("national_code = ?", nationalCode).Order("updated_at desc").First(score).Error; err != nil {
		return nil, err
	}
	return score, nil
}

func (r *postgresRepository) ReturnByID(_ context.Context, id uuid.UUID) (*Score, error) {
	score := new(Score)
	if err := r.db.Model(new(Score)).Where("id = ?", id).First(score).Error; err != nil {
		return nil, err
	}
	return score, nil
}
