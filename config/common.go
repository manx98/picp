package config

import (
	"github.com/go-ini/ini"
	"go.uber.org/zap"
	"picp/logger"
)

var Common *commonConfig

type commonConfig struct {
	cfg      *ini.Section `ini:"-"`
	LogLevel string       `ini:"log_level" validator:"required,oneof=debug info warn error dpanic panic fatal DEBUG INFO WARN ERROR DPANIC PANIC FATAL"`
	LogFile  string       `ini:"log_file" validator:"filepath"`
	BindAddr string       `ini:"bind_addr" validator:"required,hostname_port"`
	User     string       `ini:"user"`
	Password string       `ini:"password"`
}

func initCommon() {
	Common = &commonConfig{
		LogLevel: "",
		BindAddr: ":8888",
	}
	var ok bool
	Common.cfg, ok = Get("common")
	if ok {
		err := StrictMapTo(Common.cfg, Common)
		if err != nil {
			logger.Fatal("map common config failed", zap.Error(err))
		}
	}
}
