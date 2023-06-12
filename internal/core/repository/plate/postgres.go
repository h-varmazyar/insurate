package plate

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type postgresRepository struct {
	db     *gorm.DB
	logger *log.Logger
}

func NewPostgresRepository(_ context.Context, logger *log.Logger, db *gorm.DB) (Repository, error) {
	if db == nil {
		return nil, errors.New("invalid db instance for plate")
	}
	return &postgresRepository{
		db:     db,
		logger: logger,
	}, nil
}

func (r *postgresRepository) ReturnByText(ctx context.Context, text string) (*Plate, error) {
	return nil, nil
}

func (r *postgresRepository) Create(ctx context.Context, plate *Plate) error {
	return nil
}
