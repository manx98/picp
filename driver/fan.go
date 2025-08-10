package driver

import (
	"github.com/stianeikeland/go-rpio/v4"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"picp/config"
	"picp/logger"
)

const (
	maxFanHz    = 144
	maxCycleLen = 100
)

var fanPin rpio.Pin
var fanEnable atomic.Bool

func FanInit() {
	if config.Fan.Enable {
		fanPin = rpio.Pin(config.Fan.Pin)
		fanPin.Pwm()
		fanPin.Freq(maxFanHz)
		fanPin.DutyCycle(0, maxCycleLen)
	}
}

func EnableFan(enable bool) {
	if config.Fan.Enable {
		logger.Debug("change fan status", zap.Bool("enable", enable))
		if enable {
			fanEnable.Store(true)
			fanPin.DutyCycle(uint32(config.Fan.Speed), maxCycleLen)
		} else {
			fanPin.DutyCycle(0, maxCycleLen)
			fanEnable.Store(false)
		}
	} else {
		logger.Debug("fan is disabled", zap.Bool("enable", enable))
	}
}

func IsFanEnable() bool {
	return fanEnable.Load()
}
