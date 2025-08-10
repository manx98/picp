package utils

import (
	"bytes"
	"context"
	"fmt"
	"github.com/shirou/gopsutil/net"
	"go.uber.org/zap"
	"os"
	"os/exec"
	"picp/logger"
	"strconv"
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
