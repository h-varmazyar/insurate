package service

import (
	"github.com/h-varmazyar/insurate/pkg/errors"
	"net/http"
)

var (
	ErrJobNotCompleted = errors.NewWithHttp("job_not_completed", 1101, http.StatusBadRequest)
)
