package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/shirou/gopsutil/mem"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"net"
	"os/signal"
	"picp/api"
	"picp/config"
	"picp/driver"
	"picp/logger"
	"picp/utils"
	"syscall"
)

var eventChan chan driver.EventType
var apMode = false

func main() {
	flag.Parse()
	config.Init()
	driver.Init()
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	gp, ctx := errgroup.WithContext(ctx)
	eventChan = make(chan driver.EventType, 1)
	listen, err := net.Listen("tcp", config.Common.BindAddr)
	if err != nil {
		logger.Fatal("server listen error", zap.Error(err), zap.String("bind_addr", config.Common.BindAddr))
	}
	gp.Go(func() error {
		select {
		case <-ctx.Done():
			_ = listen.Close()
		}
		return nil
	})
	gp.Go(func() error {
		err := api.Run(listen)
		if err != nil {
			return fmt.Errorf("api run: %w", err)
		}
		return nil
	})
	gp.Go(func() error {
		err := driver.Run(ctx, func(event driver.EventType) error {
			select {
			case eventChan <- event:
			default:
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("driver run: %w", err)
		}
		return nil
	})
	gp.Go(func() error {
		return eventHandler(ctx)
	})
	err = gp.Wait()
	if err != nil && errors.Is(err, context.Canceled) {
		logger.Fatal("run failed", zap.Error(err))
	}
}

func eventHandler(ctx context.Context) error {
	for {
		select {
		case event := <-eventChan:
			err := eventProcessor(ctx, event)
			if err != nil {
				return fmt.Errorf("event processor: %w", err)
			}
		case <-ctx.Done():
			return fmt.Errorf("event handler: %w", context.Canceled)
		}
	}
}

func eventProcessor(ctx context.Context, event driver.EventType) error {
	if apMode {
		if event == driver.EventWifiReleased {
			apMode = false
		}
	} else {
		if event == driver.EventStatusUpdated {
			displayStatus()
		}
		if event == driver.EventWifiReleased {
			apMode = true
			displayAllAlign("WIFI Starting...")
		}
	}
	return nil
}

func displayStatus() {
	cpuTemp := driver.GetCpuTemp()
	cpuPercent := driver.GetCpuPercent()
	tx, rx := driver.GetNetSpeed()
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
	displayVerticalAlign(
		"IP "+utils.GetHostIP(),
		fmt.Sprintf("CPU %.1f%% %.2f℃", cpuPercent, cpuTemp),
		fmt.Sprintf("MEM %s %.1f%%", utils.ByteSize(memUsed, 1024), memPercent),
		fmt.Sprintf("DISK %s %.1f%%", utils.ByteSize(used, 1024), diskPercent),
		fmt.Sprintf("↑%s/s ↓%s/s", utils.ByteSize(tx, 100), utils.ByteSize(rx, 100)))
}

var statusOpt = &driver.DrawOptions{
	VerticalAlign: true,
}

func displayVerticalAlign(msg ...string) {
	err := driver.Display(statusOpt, msg...)
	if err != nil {
		logger.Warn("display vertical align error", zap.Error(err))
	}
}

var alignOpt = &driver.DrawOptions{
	HorizontalAlign: true,
	VerticalAlign:   true,
}

func displayAllAlign(msg ...string) {
	err := driver.Display(alignOpt, msg...)
	if err != nil {
		logger.Warn("display all align error", zap.Error(err))
	}
}
