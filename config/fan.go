package config

import (
	"github.com/go-ini/ini"
	"picp/logger"
)

var Fan *fanChanelCfg

type fanChanelCfg struct {
	cfg     *ini.Section `ini:"-"`
	Enable  bool         `ini:"enable"`
	Pin     int          `ini:"pin" validator:"contains=12,13,40,41,45,18,19"`
	Speed   int          `ini:"speed" validator:"gt=0,lte=100"`
	MinTemp float32      `ini:"min_temp" validator:"gte=0"`
	MaxTemp float32      `ini:"max_temp" validator:"gte=0"`
}

func newFanChanelCfg() *fanChanelCfg {
	return &fanChanelCfg{
		Enable:  false,
		Speed:   1,
		MaxTemp: 50,
		MinTemp: 45,
	}
}

func initFan() {
	Fan = newFanChanelCfg()
	var ok bool
	Fan.cfg, ok = Get("fan")
	if ok {
		if err := StrictMapTo(Fan.cfg, Fan); err != nil {
			logger.Fatalf("fan config error: %s", err)
		}
	}
}
