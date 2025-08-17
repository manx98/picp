package driver

import (
	"context"
	"github.com/stianeikeland/go-rpio/v4"
	"go.uber.org/zap"
	"picp/logger"
)

func Init(ctx context.Context) {
	err := rpio.Open()
	if err != nil {
		logger.Fatal("open rpio failed", zap.Error(err))
	}
	initStatusRunner(ctx)
	sh1106Init(ctx)
	wifiInit(ctx)
	fanInit(ctx)
}
func Close() {
	closeWifi()
	closeStatus()
	closeFan()
	closeDisplay()
}
