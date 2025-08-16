package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"net"
	"os/signal"
	"picp/api"
	"picp/config"
	"picp/driver"
	"picp/logger"
	"syscall"
)

func main() {
	flag.Parse()
	config.Init()
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	gp, ctx := errgroup.WithContext(ctx)
	driver.Init(ctx)
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
		err = api.Run(listen)
		if err != nil {
			return fmt.Errorf("api run: %w", err)
		}
		return nil
	})
	err = gp.Wait()
	if err != nil && errors.Is(err, context.Canceled) {
		logger.Fatal("run failed", zap.Error(err))
	}
}
