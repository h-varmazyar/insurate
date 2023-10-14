package service

import "github.com/h-varmazyar/insurate/pkg/amqp"

type Configs struct {
	AmqpConfigs      *amqp.Configs `json:"amqp_configs" yaml:"amqp_configs"`
	ScoreJobQueue    string        `json:"score_job_queue" yaml:"score_job_queue"`
	ScoreJobExchange string        `json:"score_job_exchange" yaml:"score_job_exchange"`
}
