package main

import (
	"context"
	"fmt"
	gormext "github.com/h-varmazyar/gopack/gorm"
	"github.com/h-varmazyar/insurate/internal/core"
	db "github.com/h-varmazyar/insurate/internal/pkg/db/PostgreSQL"
	"github.com/h-varmazyar/insurate/pkg/finnotech"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	ctx := context.Background()
	logger := log.New()

	conf, err := loadConfigs()
	if err != nil {
		log.Panic("failed to read configs")
	}

	logger.Infof("starting Insurate")

	dbInstance, err := loadDB(ctx, conf.DB)
	if err != nil {
		logger.Panicf("failed to initiate databases with error %v", err)
	}

	finnotechClient, err := finnotech.NewClient(ctx, conf.Finnotech)
	if err != nil {
		logger.WithError(err).Panicf("failed to create finnotech client")
	}

	service, err := core.NewService(ctx, logger, dbInstance, finnotechClient)
	if err != nil {
		logger.WithError(err).Panicf("failed to create service")
	}

	registerController(ctx, conf, logger, service)
}

func loadConfigs() (*Configs, error) {
	log.Infof("reding configs...")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Warnf("failed to read from env: %v", err)
		viper.AddConfigPath("./configs")  //path for docker compose configs
		viper.AddConfigPath("../configs") //path for local configs
		viper.SetConfigName("config")
		viper.SetConfigType("yml")
		if err = viper.ReadInConfig(); err != nil {
			log.Errorf("failed to read configs")
			return nil, err
		}
	}

	conf := new(Configs)
	if err := viper.Unmarshal(conf); err != nil {
		log.Errorf("faeiled unmarshal")
		return nil, err
	}

	return conf, nil
}

func loadDB(ctx context.Context, configs gormext.Configs) (*db.DB, error) {
	return db.NewDatabase(ctx, configs)
}

func registerController(_ context.Context, conf *Configs, log *log.Logger, service *core.Service) {
	controller := core.NewController(service)
	controller.RegisterRoutes()
	err := controller.Run(fmt.Sprintf(":%v", conf.HttpPort))
	if err != nil {
		log.WithError(err).Panicf("failed to start controller")
	}
}
