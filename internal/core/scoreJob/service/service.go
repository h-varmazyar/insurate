package service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	nid "github.com/h-varmazyar/gopet/national_id"
	"github.com/h-varmazyar/gopet/phone"
	"github.com/h-varmazyar/insurate/internal/core/scoreJob/repository"
	"github.com/h-varmazyar/insurate/internal/entity"
	amqpext "github.com/h-varmazyar/insurate/pkg/amqp"
	"github.com/h-varmazyar/insurate/pkg/validator"
	amqp "github.com/rabbitmq/amqp091-go"
	"strconv"
	"strings"
)

const (
	trackingIdPattern = "NSRT-%v"
)

type Service struct {
	configs            *Configs
	amqpClient         *amqpext.Client
	scoreJobRepository repository.Repository
}

func NewService(_ context.Context, configs *Configs, scoreJobRepo repository.Repository) (*Service, error) {
	amqpClient, err := amqpext.NewClient(configs.AmqpConfigs)
	if err != nil {
		return nil, err
	}

	handler := &Service{
		configs:            configs,
		amqpClient:         amqpClient,
		scoreJobRepository: scoreJobRepo,
	}

	return handler, nil
}

func (s *Service) SubmitScoreJob(ctx context.Context, submitJobReq *SubmitScoreJobRequest) (*SubmitScoreJobResponse, error) {
	if !nid.Validate(submitJobReq.NationalId) {
		return nil, ErrInvalidNationalId
	}

	details, err := phone.GetPhoneNumberDetails(submitJobReq.Mobile)
	if err != nil {
		return nil, ErrInvalidMobile
	}

	if details.Type != phone.PostPaid && details.Type != phone.PrePaid {
		return nil, ErrInvalidMobile
	}

	if !validator.IsValidPlate(submitJobReq.Plate) {
		return nil, ErrInvalidPlate
	}

	job := &entity.ScoreJob{
		NationalId:        submitJobReq.NationalId,
		Mobile:            submitJobReq.Mobile,
		Plate:             submitJobReq.Plate,
		LicenceId:         submitJobReq.LicenceId,
		InsuranceUniqueId: submitJobReq.InsuranceUniqueId,
	}

	if err = s.scoreJobRepository.Create(ctx, job); err != nil {
		return nil, err
	}

	trackingId := fmt.Sprintf(trackingIdPattern, job.ID)
	msg := amqp.Publishing{
		ContentType:  gin.MIMEJSON,
		DeliveryMode: amqp.Persistent,
		MessageId:    trackingId,
		Timestamp:    job.CreatedAt,
		Body:         []byte(job.Json()),
	}
	if err = s.amqpClient.Channel().PublishWithContext(ctx, s.configs.ScoreJobExchange, s.configs.ScoreJobQueue, false, true, msg); err != nil {
		return nil, ErrJobSubmitFailed.AddOriginalError(err)
	}
	resp := &SubmitScoreJobResponse{TrackingId: trackingId}

	return resp, nil
}

func (s *Service) JobStatus(ctx context.Context, jobStatusReq *JobStatusRequest) (*JobStatus, error) {
	strId := strings.TrimPrefix(jobStatusReq.TrackingId, "NSRT-")
	jobId, err := strconv.Atoi(strId)
	if err != nil {
		return nil, ErrInvalidJob
	}

	var status entity.JobStatus
	status, err = s.scoreJobRepository.Status(ctx, uint(jobId))
	if err != nil {
		return nil, err
	}

	resp := &JobStatus{
		Status: status.String(),
	}

	return resp, nil
}
