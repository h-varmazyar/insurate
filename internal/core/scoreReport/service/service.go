package service

import (
	"context"
	scoreJobService "github.com/h-varmazyar/insurate/internal/core/scoreJob/service"
	"github.com/h-varmazyar/insurate/internal/core/scoreReport/repository"
	"github.com/h-varmazyar/insurate/internal/entity"
	amqpext "github.com/h-varmazyar/insurate/pkg/amqp"
)

type Service struct {
	configs               *Configs
	amqpClient            *amqpext.Client
	scoreReportRepository repository.Repository
	scoreJobService       scoreJobService.Service
}

type Dependencies struct {
	ScoreReportRepository repository.Repository
	ScoreJobService       scoreJobService.Service
}

func NewService(_ context.Context, configs *Configs, dependencies *Dependencies) (*Service, error) {
	amqpClient, err := amqpext.NewClient(configs.AmqpConfigs)
	if err != nil {
		return nil, err
	}

	handler := &Service{
		configs:               configs,
		amqpClient:            amqpClient,
		scoreReportRepository: dependencies.ScoreReportRepository,
		scoreJobService:       dependencies.ScoreJobService,
	}

	return handler, nil
}

func (s *Service) ReturnReport(ctx context.Context, scoreReportReq *ScoreReportRequest) (*ScoreReport, error) {
	status, err := s.scoreJobService.JobStatus(ctx, &scoreJobService.JobStatusRequest{TrackingId: scoreReportReq.TrackingId})
	if err != nil {
		return nil, err
	}

	if status.Status != entity.JobStatusDone.String() && status.Status != entity.JobStatusFailed.String() {
		return nil, ErrJobNotCompleted
	}

	var report *entity.ScoreReport
	report, err = s.scoreReportRepository.ReturnByTrackingId(ctx, scoreReportReq.TrackingId)
	if err != nil {
		return nil, err
	}

	resp := &ScoreReport{}

	return resp, nil
}
