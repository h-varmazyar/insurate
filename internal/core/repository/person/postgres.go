package person

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
		return nil, errors.New("invalid db instance for person")
	}
	return &postgresRepository{
		db:     db,
		logger: logger,
	}, nil
}

func (r *postgresRepository) Create(ctx context.Context, person *Person) error {
	return nil
}

func (r *postgresRepository) Return(ctx context.Context, nationalCode string) (*Person, error) {
	return nil, nil
}
