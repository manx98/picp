package config

import (
	"errors"
	"picp/go-i2c"
)

type IICConfig struct {
	Enable bool `ini:"enable"`
	Bus    int  `ini:"bus" validator:"gte=0,lte=1"`
	Addr   int  `ini:"addr" validator:"get=0,lte=254"`
}

var ErrorSensorDisabled = errors.New("sensor is disabled")

func (i *IICConfig) Create() (*i2c.I2C, error) {
	if i.Enable {
		return i2c.NewI2C(uint8(i.Addr), i.Bus)
	}
	return nil, ErrorSensorDisabled
}
