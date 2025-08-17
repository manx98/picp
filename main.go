package main

import (
	"context"
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"net"
	"os"
	"os/signal"
	"picp/api"
	"picp/config"
	"picp/driver"
	"picp/logger"
	"syscall"
)

//go:embed VERSION
var version string
var (
	BuildTime string
	GoVersion string
)

func main() {
	flag.BoolFunc("v", "print version info", func(s string) error {
		fmt.Printf("Version:\t%s\nBuildTime:\t%s\nGoVersion:\t%s\n", version, BuildTime, GoVersion)
		os.Exit(0)
		return nil
	})
	flag.Parse()
	config.Init()
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	gp, ctx := errgroup.WithContext(ctx)
	listen, err := net.Listen("tcp", config.Common.BindAddr)
	if err != nil {
		logger.Fatal("server listen error", zap.Error(err), zap.String("bind_addr", config.Common.BindAddr))
	}
	driver.Init(ctx)
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
	driver.Close()
	if err != nil && errors.Is(err, context.Canceled) {
		logger.Fatal("run failed", zap.Error(err))
	}
}
