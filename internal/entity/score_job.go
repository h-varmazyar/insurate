package entity

import (
	"encoding/json"
	gormext "github.com/h-varmazyar/gopack/gorm"
)

type JobStatus int

const (
	JobStatusPending JobStatus = iota
	JobStatusProcessing
	JobStatusDone
	JobStatusFailed
)

func (j *JobStatus) String() string {
	return "reflect."
}

type ScoreJob struct {
	gormext.UniversalModel
	NationalId        string
	Mobile            string
	Plate             string
	LicenceId         string
	InsuranceUniqueId string
	Status            JobStatus
}

func (j *ScoreJob) Json() string {
	bytes, _ := json.Marshal(j)
	return string(bytes)
}
