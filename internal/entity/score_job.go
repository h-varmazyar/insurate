package entity

import (
	"encoding/json"
	"gorm.io/gorm"
)

type JobStatus int

const (
	JobStatusUnknown JobStatus = iota
	JobStatusPending
	JobStatusProcessing
	JobStatusDone
	JobStatusFailed
)

var (
	jobStatusName = map[JobStatus]string{
		JobStatusUnknown:    "Unknown",
		JobStatusPending:    "Pending",
		JobStatusProcessing: "Processing",
		JobStatusDone:       "Done",
		JobStatusFailed:     "Failed",
	}

	jobStatusValue = map[string]JobStatus{
		"Unknown":    JobStatusUnknown,
		"Pending":    JobStatusPending,
		"Processing": JobStatusProcessing,
		"Done":       JobStatusDone,
		"Failed":     JobStatusFailed,
	}
)

func (j JobStatus) String() string {
	name, ok := jobStatusName[j]
	if ok {
		return name
	}
	return ""
}

func (j *JobStatus) Value(value interface{}) error {
	*j = jobStatusValue[value.(string)]
	return nil
}

type ScoreJob struct {
	gorm.Model
	TrackingId        string
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
