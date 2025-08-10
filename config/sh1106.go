package config

import (
	"github.com/go-ini/ini"
	"go.uber.org/zap"
	"picp/logger"
	"picp/sh1106"
)

var SH1106 *sh1106Config

type sh1106Config struct {
	cfg       *ini.Section `ini:"-"`
	IICConfig `ini:",extends"`
	Height    int `ini:"height" validator:"gt=0,lt=32767"`
	Width     int `ini:"width" validator:"gt=0,lt=32767"`
	VccState  int `ini:"vcc_state" validator:"get=0,lte=1"`
}

func (c *sh1106Config) GetMode() sh1106.VccMode {
	return []sh1106.VccMode{
		sh1106.ExternalVCC,
		sh1106.SwitchCAPVCC,
	}[c.VccState]
}

func newSh1106Config() *sh1106Config {
	return &sh1106Config{
		IICConfig: IICConfig{
			Addr: 0x3C,
			Bus:  1,
		},
		Height: 64,
		Width:  128,
	}
}

func initSH1106() {
	SH1106 = newSh1106Config()
	var ok bool
	SH1106.cfg, ok = Get("sh1106")
	if ok {
		err := StrictMapTo(SH1106.cfg, SH1106)
		if err != nil {
			logger.Fatal("failed to load sh1106 config", zap.Error(err))
		}
	}
}
