package db

import gormext "github.com/h-varmazyar/gopack/gorm"

type Migration struct {
	gormext.UniversalModel
	TableName   string
	Tag         string
	Description string
}

const MigrationTable = "migrations"
