package plate

import (
	"context"
	gormext "github.com/h-varmazyar/gopack/gorm"
	personRepo "github.com/h-varmazyar/insurate/internal/core/repository/person"
	db "github.com/h-varmazyar/insurate/internal/pkg/db/PostgreSQL"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const tableName = "plates"

type Plate struct {
	gormext.UniversalModel
	Person      *personRepo.Person
	Alphabet    string
	StartNumber int8
	EndNumber   int8
	RegionCode  int8
	Text        string
}

type Repository interface {
	ReturnByText(ctx context.Context, text string) (*Plate, error)
	Create(ctx context.Context, plate *Plate) error
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
			err = tx.AutoMigrate(new(Plate))
			if err != nil {
				return err
			}
			newMigrations = append(newMigrations, &db.Migration{
				TableName:   tableName,
				Tag:         "v1.0.0",
				Description: "create plates table",
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
