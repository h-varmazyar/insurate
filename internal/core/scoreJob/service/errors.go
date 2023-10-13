package service

import (
	"github.com/h-varmazyar/insurate/pkg/errors"
	"net/http"
)

var (
	ErrInvalidNationalId = errors.NewWithHttp("invalid_national_id", 1001, http.StatusBadRequest)
	ErrInvalidMobile     = errors.NewWithHttp("invalid_mobile", 1002, http.StatusBadRequest)
	ErrInvalidPlate      = errors.NewWithHttp("invalid_plate", 1003, http.StatusBadRequest)
	ErrInvalidJob        = errors.NewWithHttp("invalid_job", 1004, http.StatusBadRequest)
	ErrJobNotPrepared    = errors.NewWithHttp("job_not_prepared", 1005, http.StatusTooEarly)
	ErrJobSubmitFailed   = errors.NewWithHttp("job_submit_failed", 1006, http.StatusBadRequest)
)
