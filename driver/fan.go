package driver

import (
	"context"
	"github.com/stianeikeland/go-rpio/v4"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"picp/config"
	"picp/logger"
	"picp/utils"
	"time"
)

const (
	maxFanHz    = 144
	maxCycleLen = 100
)

var fanEnable atomic.Bool

var fanRunner *utils.Runner

func fanInit(ctx context.Context) {
	fanRunner = utils.NewRunner(ctx, runFan)
	fanRunner.Start()
}

func runFan(ctx context.Context) {
	cfg := config.GetFanCfg()
	if !cfg.Enable {
		return
	}
	fanPin := rpio.Pin(cfg.Pin)
	fanPin.Pwm()
	fanPin.Freq(maxFanHz)
	fanPin.DutyCycle(0, maxCycleLen)
	tik := time.NewTicker(3 * time.Second)
	defer func() {
		tik.Stop()
		changeFanSpeed(fanPin, 0)
	}()
	for ctx.Err() == nil {
		temperature, err := utils.GetCpuTemperature()
		if err != nil {
			logger.Debug("cpu temperature error", zap.Error(err))
		} else {
			if temperature > cfg.MaxTemp && !fanEnable.Load() {
				changeFanSpeed(fanPin, uint32(cfg.Speed))
			} else if temperature < cfg.MinTemp && fanEnable.Load() {
				changeFanSpeed(fanPin, 0)
			}
		}
		select {
		case <-tik.C:
		case <-ctx.Done():
			return
		}
	}
}

func changeFanSpeed(fanPin rpio.Pin, speed uint32) {
	logger.Debug("change fan speed", zap.Uint32("speed", speed))
	fanEnable.Store(speed != 0)
	fanPin.DutyCycle(speed, maxCycleLen)
}

func SetFanConfig(cfg *config.FanChanelCfg) error {
	err := config.SetFanCfg(cfg)
	if err != nil {
		return err
	}
	_ = fanRunner.Stop(context.Background())
	fanRunner.Start()
	return nil
}

func closeFan() {
	_ = fanRunner.Stop(context.Background())
}
