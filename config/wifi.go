package config

import (
	"github.com/go-ini/ini"
	"picp/logger"
)

var WIFI *wifiConfig

type wifiConfig struct {
	cfg      *ini.Section `ini:"-"`
	Enable   bool         `ini:"enable"`
	Pin      int          `ini:"pin" validator:"gte=0,lt=255"`
	SSID     string       `ini:"ssid" validator:"max=32"`
	Password string       `ini:"password" validator:"min=8"`
}

func initWifi() {
	WIFI = new(wifiConfig)
	var ok bool
	WIFI.cfg, ok = Get("wifi")
	if ok {
		err := StrictMapTo(WIFI.cfg, WIFI)
		if err != nil {
			logger.Fatalf("ap config error: %s", err)
		}
	}
}
