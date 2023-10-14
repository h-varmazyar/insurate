package scoreReport

import (
	"context"
	scoreJobService "github.com/h-varmazyar/insurate/internal/core/scoreJob/service"
	"github.com/h-varmazyar/insurate/internal/core/scoreReport/repository"
	"github.com/h-varmazyar/insurate/internal/core/scoreReport/service"
	db "github.com/h-varmazyar/insurate/pkg/db/PostgreSQL"
	log "github.com/sirupsen/logrus"
)

type App struct {
	service    *service.Service
	repository repository.Repository
}

type Dependencies struct {
	DB              *db.DB
	ScoreJobService scoreJobService.Service
}

func NewApp(ctx context.Context, logger *log.Logger, configs *Configs, dependencies *Dependencies) (*App, error) {
	var err error
	app := new(App)
	app.repository, err = repository.NewRepository(ctx, logger, dependencies.DB)
	if err != nil {
		return nil, err
	}

	serviceDependencies := &service.Dependencies{
		ScoreReportRepository: app.repository,
		ScoreJobService:       dependencies.ScoreJobService,
	}

	app.service, err = service.NewService(ctx, configs.ServiceConfigs, serviceDependencies)
	if err != nil {
		return nil, err
	}

	return app, nil
}
