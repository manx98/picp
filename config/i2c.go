package config

import (
	"errors"
	"picp/go-i2c"
)

type IICConfig struct {
	Enable bool `json:"enable" ini:"enable"`
	Bus    int  `json:"bus" ini:"bus,omitempty" validate:"required,gte=0,lte=1"`
	Addr   int  `json:"addr" ini:"addr,omitempty" validate:"required,gte=0,lte=254"`
}

var ErrorSensorDisabled = errors.New("sensor is disabled")

func (i *IICConfig) Create() (*i2c.I2C, error) {
	if i.Enable {
		return i2c.NewI2C(uint8(i.Addr), i.Bus)
	}
	return nil, ErrorSensorDisabled
}
