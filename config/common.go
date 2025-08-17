package config

import (
	"github.com/go-ini/ini"
	"go.uber.org/zap"
	"picp/logger"
	"sync"
)

var Common = common{
	LogLevel:     "info",
	BindAddr:     ":8888",
	CookieMaxAge: 168,
}
var cfgLock sync.RWMutex

type common struct {
	cfg          *ini.Section `ini:"-"`
	LogLevel     string       `ini:"log_level" validate:"required,oneof=debug info warn error dpanic panic fatal DEBUG INFO WARN ERROR DPANIC PANIC FATAL"`
	LogFile      string       `ini:"log_file"`
	BindAddr     string       `ini:"bind_addr" validate:"required,hostname_port"`
	User         string       `ini:"user"`
	Password     string       `ini:"password"`
	CookieMaxAge int          `ini:"cookie_max_age" validate:"gt=0"`
}

func initCommon() {
	var ok bool
	Common.cfg, ok = Get("common")
	if ok {
		err := StrictMapTo(Common.cfg, &Common)
		if err != nil {
			logger.Fatal("map common config failed", zap.Error(err))
		}
	}
	Common.CookieMaxAge *= 60 * 60
}
func SaveCfg() error {
	return rootCfg.SaveTo(*cfgPath)
}

func SetLoginSetting(username string, password string, maxAge int) (err error) {
	old := Common
	Common.User = username
	Common.Password = password
	if maxAge > 0 {
		Common.CookieMaxAge = maxAge
	}
	defer func() {
		if err != nil {
			Common = old
		}
	}()
	err = Common.cfg.ReflectFrom(&Common)
	if err != nil {
		return err
	}
	return SaveCfg()
}
