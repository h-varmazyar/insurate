package scoreJob

import (
	"context"
	"github.com/h-varmazyar/insurate/internal/core/scoreJob/repository"
	"github.com/h-varmazyar/insurate/internal/core/scoreJob/service"
	db "github.com/h-varmazyar/insurate/pkg/db/PostgreSQL"
	log "github.com/sirupsen/logrus"
)

type App struct {
	service    *service.Service
	repository repository.Repository
}

func NewApp(ctx context.Context, logger *log.Logger, configs *Configs, db *db.DB) (*App, error) {
	var err error
	app := new(App)
	app.repository, err = repository.NewRepository(ctx, logger, db)
	if err != nil {
		return nil, err
	}

	app.service, err = service.NewService(ctx, configs.ServiceConfigs, app.repository)
	if err != nil {
		return nil, err
	}

	return app, nil
}
