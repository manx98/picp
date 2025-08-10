package driver

import (
	"context"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/stianeikeland/go-rpio/v4"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"picp/config"
	"picp/logger"
	"picp/utils"
	"time"
)

var cpuPercent atomic.Float64
var cpuTemp atomic.Float64
var lastTxCount uint64
var lastRxCount uint64
var lastCountTime time.Time
var txSpeed atomic.Int64
var rxSpeed atomic.Int64

type EventType uint8

const (
	EventWifiPressed EventType = iota
	EventWifiReleased
	EventStatusUpdated
)

type EventHandler func(event EventType) error

func Init() {
	err := rpio.Open()
	if err != nil {
		logger.Fatal("open rpio failed", zap.Error(err))
	}
	FanInit()
	sh1106Init()
	wifiInit()
}

func Run(ctx context.Context, handler EventHandler) error {
	gp, gpCtx := errgroup.WithContext(ctx)
	gp.Go(func() error {
		timer := time.NewTicker(time.Second)
		defer timer.Stop()
		for {
			select {
			case <-timer.C:
				updateStatus()
				err := handler(EventStatusUpdated)
				if err != nil {
					return fmt.Errorf("status event handler: %w", err)
				}
			case <-gpCtx.Done():
				return fmt.Errorf("sensor update: %w", context.Canceled)
			}
		}
	})
	if config.WIFI.Enable {
		gp.Go(func() error {
			timer := time.NewTicker(time.Millisecond * 100)
			defer timer.Stop()
			lastApPress := false
			for {
				select {
				case <-timer.C:
					var err error
					if isWifi() != lastApPress {
						if lastApPress {
							err = handler(EventWifiReleased)
						} else {
							err = handler(EventWifiPressed)
						}
						if err != nil {
							return fmt.Errorf("ap button handler: %w", err)
						}
						lastApPress = !lastApPress
					}
				case <-gpCtx.Done():
					return fmt.Errorf("WIFI button scan: %w", context.Canceled)
				}
			}
		})
	}
	return gp.Wait()
}

func updateStatus() {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		logger.Debug("cpu percent error", zap.Error(err))
	} else if len(percent) > 0 {
		cpuPercent.Store(percent[0])
	}
	if config.Fan.Enable {
		temperature, err := utils.GetCpuTemperature()
		if err != nil {
			cpuTemp.Store(0)
			logger.Debug("cpu temperature error", zap.Error(err))
		} else {
			cpuTemp.Store(float64(temperature))
			if temperature > config.Fan.MaxTemp && !IsFanEnable() {
				EnableFan(true)
			} else if temperature < config.Fan.MinTemp && IsFanEnable() {
				EnableFan(false)
			}
		}
	}
	updateSpeedCount()
}

func GetCpuPercent() float64 {
	return cpuPercent.Load()
}

func GetCpuTemp() float64 {
	return cpuTemp.Load()
}

func updateSpeedCount() {
	txCount, rxCount, err := utils.GetNetIoCounters(context.Background())
	if err != nil {
		logger.Debug("get net io counters error", zap.Error(err))
		return
	}
	current := time.Now()
	if !lastCountTime.IsZero() {
		duration := time.Now().Sub(lastCountTime).Seconds()
		if txCount > lastTxCount {
			txSpeed.Store(int64(float64(txCount-lastTxCount) / duration))
		}
		if rxCount > lastRxCount {
			rxSpeed.Store(int64(float64(rxCount-lastRxCount) / duration))
		}
	}
	lastTxCount = txCount
	lastRxCount = rxCount
	lastCountTime = current
}

func GetNetSpeed() (tx, rx int64) {
	return txSpeed.Load(), rxSpeed.Load()
}
