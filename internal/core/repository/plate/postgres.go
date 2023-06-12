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

func (r *postgresRepository) ReturnByText(_ context.Context, text string) (*Plate, error) {
	plate := new(Plate)
	if err := r.db.Model(new(Plate)).Where("text = ?", text).Error; err != nil {
		return nil, err
	}
	return plate, nil
}

func (r *postgresRepository) Create(_ context.Context, plate *Plate) error {
	if err := r.db.Save(plate).Error; err != nil {
		return err
	}
	return nil
}
