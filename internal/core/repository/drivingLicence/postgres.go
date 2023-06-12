package drivingLicence

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
		return nil, errors.New("invalid db instance for driving licence")
	}
	return &postgresRepository{
		db:     db,
		logger: logger,
	}, nil
}

func (r *postgresRepository) ReturnByNumber(_ context.Context, number uint64) (*DrivingLicence, error) {
	licence := new(DrivingLicence)
	if err := r.db.Model(new(DrivingLicence)).Where("number = ?", number).First(licence).Error; err != nil {
		return nil, err
	}
	return licence, nil
}

func (r *postgresRepository) Create(_ context.Context, licence *DrivingLicence) error {
	if err := r.db.Save(licence).Error; err != nil {
		return err
	}
	return nil
}
