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

func (r *postgresRepository) Create(ctx context.Context, score *Score) error {
	return nil
}

func (r *postgresRepository) Update(ctx context.Context, score *Score) error {
	return nil
}

func (r *postgresRepository) UpdateStatus(ctx context.Context, scoreID uuid.UUID, newStatus Status) error {
	return nil
}

func (r *postgresRepository) Last(ctx context.Context, nationalCode string) (*Score, error) {
	return nil, nil
}

func (r *postgresRepository) ReturnByID(ctx context.Context, id uuid.UUID) (*Score, error) {
	return nil, nil
}
