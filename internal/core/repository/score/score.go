package score

import (
	"context"
	"github.com/google/uuid"
	gormext "github.com/h-varmazyar/gopack/gorm"
	drivingLicenceRepo "github.com/h-varmazyar/insurate/internal/core/repository/drivingLicence"
	"github.com/h-varmazyar/insurate/internal/core/repository/drivingOffence"
	personRepo "github.com/h-varmazyar/insurate/internal/core/repository/person"
	plateRepo "github.com/h-varmazyar/insurate/internal/core/repository/plate"
	db "github.com/h-varmazyar/insurate/internal/pkg/db/PostgreSQL"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const tableName = "scores"

type Status int8

const (
	Pending Status = iota
	PreparingData
	Calculating
	Done
	Failed
)

type Score struct {
	gormext.UniversalModel
	NationalCode    string
	Person          *personRepo.Person                 `gorm:"-"`
	DrivingLicence  *drivingLicenceRepo.DrivingLicence `gorm:"-"`
	DrivingOffences []*drivingOffence.DrivingOffence   `gorm:"-"`
	Plate           *plateRepo.Plate                   `gorm:"-"`
	NegativeScore   int8                               `gorm:"-"`
	Value           float64
	Status          Status
}

type Repository interface {
	Create(ctx context.Context, score *Score) error
	Update(ctx context.Context, score *Score) error
	UpdateStatus(ctx context.Context, scoreID uuid.UUID, newStatus Status) error
	Last(ctx context.Context, nationalCode string) (*Score, error)
	ReturnByID(ctx context.Context, id uuid.UUID) (*Score, error)
}

func NewRepository(ctx context.Context, logger *log.Logger, db *db.DB) (Repository, error) {
	if err := migration(ctx, db); err != nil {
		return nil, err
	}
	return NewPostgresRepository(ctx, logger, db.PostgresDB)
}

func migration(_ context.Context, dbInstance *db.DB) error {
	var err error
	migrations := make(map[string]struct{})
	tags := make([]string, 0)
	err = dbInstance.PostgresDB.Table(db.MigrationTable).Where("table_name = ?", tableName).Select("tag").Find(&tags).Error
	if err != nil {
		return err
	}

	for _, tag := range tags {
		migrations[tag] = struct{}{}
	}

	newMigrations := make([]*db.Migration, 0)
	err = dbInstance.PostgresDB.Transaction(func(tx *gorm.DB) error {
		if _, ok := migrations["v1.0.0"]; !ok {
			err = tx.AutoMigrate(new(Score))
			if err != nil {
				return err
			}
			newMigrations = append(newMigrations, &db.Migration{
				TableName:   tableName,
				Tag:         "v1.0.0",
				Description: "create scores table",
			})
		}
		err = tx.Model(new(db.Migration)).CreateInBatches(&newMigrations, 100).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
