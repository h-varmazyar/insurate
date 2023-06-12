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

func (r *postgresRepository) Create(_ context.Context, person *Person) error {
	if err := r.db.Save(person).Error; err != nil {
		return err
	}
	return nil
}

func (r *postgresRepository) Return(_ context.Context, nationalCode string) (*Person, error) {
	person := new(Person)
	if err := r.db.Model(new(Person)).Where("national_code = ?", nationalCode).Error; err != nil {
		return nil, err
	}
	return person, nil
}
