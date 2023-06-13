package main

import (
	gormext "github.com/h-varmazyar/gopack/gorm"
	"github.com/h-varmazyar/insurate/pkg/finnotech"
)

type Configs struct {
	HttpPort  uint16            `mapstructure:"http_port"`
	Finnotech *finnotech.Config `mapstructure:"finnotech"`
	DB        gormext.Configs   `mapstructure:"db"`
}
