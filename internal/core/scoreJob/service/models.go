package service

type SubmitScoreJobRequest struct {
	Mobile            string
	NationalId        string
	LicenceId         string
	InsuranceUniqueId string
	Plate             string
}

type SubmitScoreJobResponse struct {
	TrackingId string
}

type JobStatusRequest struct {
	TrackingId string
}

type JobStatus struct {
	Status string
}
