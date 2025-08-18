package driver

import (
	"context"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"picp/config"
	"picp/logger"
	"picp/utils"
	"sync"
	"time"
)

var statusRunner StatusRunner

type StatusRunner struct {
	*utils.Runner
	cpuPercent    atomic.Float64
	lastTxCount   uint64
	lastRxCount   uint64
	lastCountTime time.Time
	txSpeed       atomic.Int64
	rxSpeed       atomic.Int64
	statusLock    sync.Mutex
	statusEnabled atomic.Bool
	lastMsg       []string
}

func initStatusRunner(ctx context.Context) {
	statusRunner.statusEnabled.Store(true)
	statusRunner.Runner = utils.NewRunner(ctx, statusRunner.statusHandler)
}

func (s *StatusRunner) statusHandler(ctx context.Context) {
	interval := time.Second * time.Duration(config.SH1106.StatusInterval)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for ctx.Err() == nil {
		if s.statusEnabled.Load() {
			s.updateStatus(ctx, interval)
		}
		select {
		case <-ticker.C:
		case <-ctx.Done():
			return
		}
	}
}

func (s *StatusRunner) updateStatus(ctx context.Context, interval time.Duration) {
	percent, err := cpu.PercentWithContext(ctx, interval, false)
	if err != nil {
		s.cpuPercent.Store(-1)
		logger.Debug("cpu percent error", zap.Error(err))
		return
	} else if len(percent) > 0 {
		s.cpuPercent.Store(percent[0])
	}
	s.updateSpeedCount(ctx)
	s.statusLock.Lock()
	defer s.statusLock.Unlock()
	if s.statusEnabled.Load() {
		s.DisplayStatus()
	}
}

func (s *StatusRunner) updateSpeedCount(ctx context.Context) {
	txCount, rxCount, err := utils.GetNetIoCounters(ctx)
	if err != nil {
		logger.Debug("get net io counters error", zap.Error(err))
		s.lastCountTime = time.Time{}
		s.txSpeed.Store(-1)
		s.rxSpeed.Store(-1)
		return
	}
	current := time.Now()
	if !s.lastCountTime.IsZero() {
		duration := time.Now().Sub(s.lastCountTime).Seconds()
		if txCount > s.lastTxCount {
			s.txSpeed.Store(int64(float64(txCount-s.lastTxCount) / duration))
		}
		if rxCount > s.lastRxCount {
			s.txSpeed.Store(int64(float64(rxCount-s.lastRxCount) / duration))
		}
	}
	s.lastTxCount = txCount
	s.lastRxCount = rxCount
	s.lastCountTime = current
}

func (s *StatusRunner) StatusShowEnable(enabled bool) {
	s.statusLock.Lock()
	defer s.statusLock.Unlock()
	s.statusEnabled.Store(enabled)
	if enabled {
		s.DisplayStatus()
	}
}

func (s *StatusRunner) DisplayStatus() {
	cpuTemp, err := utils.GetCpuTemperature()
	if err != nil {
		cpuTemp = -1
		logger.Debug("cpu temperature error", zap.Error(err))
	}
	total, used, err := utils.GetRootDiskInfo()
	var diskPercent float64
	if err != nil {
		logger.Warn("get root disk info error", zap.Error(err))
	}
	if used > 0 {
		diskPercent = float64(used) / float64(total) * 100
	}
	var memTotal, memUsed int64
	var memPercent float64
	memory, err := mem.VirtualMemory()
	if err == nil {
		memTotal = int64(memory.Total)
		memUsed = int64(memory.Used)
		if memTotal > 0 {
			memPercent = float64(memUsed) / float64(memTotal) * 100
		}
	} else {
		logger.Warn("get memory info error", zap.Error(err))
	}
	DisplayVerticalAlign("IP "+utils.GetHostIP(),
		fmt.Sprintf("CPU %.1f%% %.2f℃", s.cpuPercent.Load(), cpuTemp),
		fmt.Sprintf("MEM %s %.1f%%", utils.ByteSize(memUsed, 1024), memPercent),
		fmt.Sprintf("DISK %s %.1f%%", utils.ByteSize(used, 1024), diskPercent),
		fmt.Sprintf("↑%s/s ↓%s/s", utils.ByteSize(s.txSpeed.Load(), 100), utils.ByteSize(s.rxSpeed.Load(), 100)))
}

func closeStatus() {
	_ = statusRunner.Stop(context.Background())
}
