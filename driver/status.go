package driver

import (
	"context"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"picp/logger"
	"picp/utils"
	"sync"
	"time"
)

var statusRunner *utils.Runner
var cpuPercent float64
var cpuTemp float32
var lastTxCount uint64
var lastRxCount uint64
var lastCountTime time.Time
var txSpeed int64
var rxSpeed int64
var statusLock sync.Mutex
var statusEnabled atomic.Bool
var lastMsg []string

func initStatusRunner(ctx context.Context) {
	statusEnabled.Store(true)
	statusRunner = utils.NewRunner(ctx, statusHandler)
}

func statusHandler(ctx context.Context) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	for ctx.Err() == nil {
		if statusEnabled.Load() {
			updateStatus(ctx)
		}
		select {
		case <-ticker.C:
		case <-ctx.Done():
			return
		}
	}
}

func updateStatus(ctx context.Context) {
	percent, err := cpu.PercentWithContext(ctx, time.Second, false)
	if err != nil {
		cpuPercent = -1
		logger.Debug("cpu percent error", zap.Error(err))
		return
	} else if len(percent) > 0 {
		cpuPercent = percent[0]
	}
	cpuTemp, err = utils.GetCpuTemperature()
	if err != nil {
		cpuTemp = -1
		logger.Debug("cpu temperature error", zap.Error(err))
	}
	updateSpeedCount(ctx)
	displayStatus()
}

func updateSpeedCount(ctx context.Context) {
	txCount, rxCount, err := utils.GetNetIoCounters(ctx)
	if err != nil {
		logger.Debug("get net io counters error", zap.Error(err))
		lastCountTime = time.Time{}
		txSpeed = -1
		rxSpeed = -1
		return
	}
	current := time.Now()
	if !lastCountTime.IsZero() {
		duration := time.Now().Sub(lastCountTime).Seconds()
		if txCount > lastTxCount {
			txSpeed = int64(float64(txCount-lastTxCount) / duration)
		}
		if rxCount > lastRxCount {
			rxSpeed = int64(float64(rxCount-lastRxCount) / duration)
		}
	}
	lastTxCount = txCount
	lastRxCount = rxCount
	lastCountTime = current
}

func StatusShowEnable(enabled bool) {
	statusLock.Lock()
	defer statusLock.Unlock()
	statusEnabled.Store(enabled)
	DisplayAllAlign(lastMsg...)
}

func displayStatus() {
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
	statusLock.Lock()
	defer statusLock.Unlock()
	if statusEnabled.Load() {
		lastMsg = []string{
			"IP " + utils.GetHostIP(),
			fmt.Sprintf("CPU %.1f%% %.2f℃", cpuPercent, cpuTemp),
			fmt.Sprintf("MEM %s %.1f%%", utils.ByteSize(memUsed, 1024), memPercent),
			fmt.Sprintf("DISK %s %.1f%%", utils.ByteSize(used, 1024), diskPercent),
			fmt.Sprintf("↑%s/s ↓%s/s", utils.ByteSize(txSpeed, 100), utils.ByteSize(rxSpeed, 100)),
		}
		DisplayVerticalAlign(lastMsg...)
	}
}
