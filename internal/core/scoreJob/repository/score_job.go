package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/h-varmazyar/insurate/internal/entity"
	db "github.com/h-varmazyar/insurate/pkg/db/PostgreSQL"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	tableName = "score_jobs"
)

type Repository interface {
	Create(ctx context.Context, job *entity.ScoreJob) error
	Status(ctx context.Context, jobId uuid.UUID) (entity.JobStatus, error)
}

func NewRepository(ctx context.Context, logger *log.Logger, db *db.DB) (Repository, error) {
	if err := migration(ctx, db); err != nil {
		return nil, err
	}
	return newPostgresRepository(ctx, logger, db.PostgresDB)
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
			err = tx.AutoMigrate(new(entity.ScoreJob))
			if err != nil {
				return err
			}
			newMigrations = append(newMigrations, &db.Migration{
				TableName:   tableName,
				Tag:         "v1.0.0",
				Description: "create score_jobs table",
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
