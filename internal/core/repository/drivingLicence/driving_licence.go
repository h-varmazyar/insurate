package drivingLicence

import (
	"context"
	gormext "github.com/h-varmazyar/gopack/gorm"
	personRepo "github.com/h-varmazyar/insurate/internal/core/repository/person"
	db "github.com/h-varmazyar/insurate/internal/pkg/db/PostgreSQL"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

const tableName = "driving_licences"

type LicenceType int8
type LicenceStatus int8

const (
	LicenceMotorCycle LicenceType = iota
	LicenceMotorCycle200
	LicenceNormalCar
	LicenceMiddleCar
	LicenceLargeCar
)

const (
	LicenceStatusAllowed LicenceStatus = iota
	LicenceStatusNotAllowed
)

type DrivingLicence struct {
	gormext.UniversalModel
	Person         *personRepo.Person
	Number         uint64 //primary key
	ExpirationTime time.Time
	AllowedVehicle []AllowedLicence
	NegativeScore  int8
	OffenceCount   int16
	Rule           string
	IssuedDate     time.Time
}

type AllowedLicence struct {
	Type       LicenceType
	Status     LicenceStatus
	IssuedDate time.Time
}

type Repository interface {
	ReturnByNumber(ctx context.Context, number uint64) (*DrivingLicence, error)
	Create(ctx context.Context, licence *DrivingLicence) error
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
			err = tx.AutoMigrate(new(DrivingLicence))
			if err != nil {
				return err
			}
			newMigrations = append(newMigrations, &db.Migration{
				TableName:   tableName,
				Tag:         "v1.0.0",
				Description: "create driving_licences table",
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
