package utils

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/shirou/gopsutil/net"
	"go.uber.org/zap"
	"os"
	"os/exec"
	"picp/logger"
	"strconv"
	"sync"
	"syscall"
)

func GetHostIP() string {
	output, err := exec.Command("/bin/bash", "-c", "hostname -I | cut -d ' ' -f1").CombinedOutput()
	if err != nil {
		logger.Warn("failed to get host ip", zap.Error(err))
		return ""
	}
	return trimNumber(output)
}

func trimNumber(value []byte) string {
	return string(bytes.Trim(value, " \t\n\r\v\f "))
}

func GetCpuTemperature() (float32, error) {
	file, err := os.ReadFile("/sys/class/thermal/thermal_zone0/temp")
	if err != nil {
		return 0, fmt.Errorf("read cpu temp file: %w", err)
	}
	temp, err := strconv.Atoi(trimNumber(file))
	if err != nil {
		return 0, fmt.Errorf("convert cpu temp value: %w", err)
	}
	return float32(temp) / 1000, nil
}

func GetRootDiskInfo() (total, used int64, err error) {
	var stat syscall.Statfs_t
	err = syscall.Statfs("/", &stat)
	if err != nil {
		return
	}
	total = int64(stat.Blocks) * stat.Bsize
	used = total - int64(stat.Bavail)*stat.Bsize
	return
}

func GetNetIoCounters(ctx context.Context) (up, down uint64, err error) {
	counters, err := net.IOCountersWithContext(ctx, true)
	if err != nil {
		return 0, 0, err
	}
	for _, counter := range counters {
		up += counter.BytesSent
		down += counter.BytesRecv
	}
	return
}

func ByteSize(bytes int64, minValue float64) string {
	tmp := float64(bytes)
	for _, unit := range []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"} {
		if tmp < minValue {
			// 保证格式化后的长度不超过 5 个字符
			if tmp < 10 {
				return fmt.Sprintf("%.2f%s", tmp, unit)
			} else {
				return fmt.Sprintf("%.1f%s", tmp, unit)
			}
		} else {
			tmp /= 1024
		}
	}
	return fmt.Sprintf("%dYB", bytes)
}

type Runner struct {
	callback func(ctx context.Context)
	rootCtx  context.Context
	ctx      context.Context
	stop     context.CancelFunc
	doneCtx  context.Context
	done     context.CancelFunc
	lock     sync.Mutex
}

func (tr *Runner) run() {
	defer func() {
		tr.lock.Lock()
		defer tr.lock.Unlock()
		tr.done()
		tr.done = nil
	}()
	tr.callback(tr.ctx)
}

func (tr *Runner) Start() {
	tr.lock.Lock()
	defer tr.lock.Unlock()
	if tr.done == nil {
		tr.ctx, tr.stop = context.WithCancel(tr.rootCtx)
		tr.doneCtx, tr.done = context.WithCancel(context.Background())
		go tr.run()
	}
}

func (tr *Runner) Stop(ctx context.Context) error {
	tr.lock.Lock()
	needStop := tr.done != nil
	doneCtx := tr.doneCtx
	if needStop {
		tr.stop()
	}
	tr.lock.Unlock()
	if needStop {
		select {
		case <-doneCtx.Done():
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

func NewRunner(ctx context.Context, callback func(ctx context.Context)) *Runner {
	tr := &Runner{
		rootCtx:  ctx,
		callback: callback,
	}
	return tr
}

func Sha1Sum(data string) string {
	sum := sha1.Sum([]byte(data))
	return hex.EncodeToString(sum[:])
}
