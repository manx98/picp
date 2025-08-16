package config

import (
	"github.com/go-ini/ini"
	"go.uber.org/zap"
	"picp/logger"
	"picp/sh1106"
	"sync"
)

var SH1106 = SH1106Config{
	IICConfig: IICConfig{
		Enable: true,
		Bus:    1,
		Addr:   0x3C,
	},
	Width:    128,
	Height:   64,
	VccState: 0,
}
var sh1106Lock sync.Mutex

type SH1106Config struct {
	cfg       *ini.Section `ini:"-"`
	IICConfig `ini:",extends"`
	Height    int `json:"height" ini:"height,omitempty" validator:"gt=0,lt=32767"`
	Width     int `json:"width" ini:"width,omitempty" validator:"gt=0,lt=32767"`
	VccState  int `json:"vcc_state" ini:"vcc_state,omitempty" validator:"get=0,lte=1"`
}

func (c *SH1106Config) GetMode() sh1106.VccMode {
	return []sh1106.VccMode{
		sh1106.ExternalVCC,
		sh1106.SwitchCAPVCC,
	}[c.VccState]
}

func initSH1106() {
	var ok bool
	SH1106.cfg, ok = Get("sh1106")
	if ok {
		err := StrictMapTo(SH1106.cfg, &SH1106)
		if err != nil {
			logger.Fatal("failed to load sh1106 config", zap.Error(err))
		}
	}
}

func GetSH1106Cfg() SH1106Config {
	sh1106Lock.Lock()
	defer sh1106Lock.Unlock()
	return SH1106
}

func SaveSH1106(cfg *SH1106Config) error {
	SH1106.IICConfig = cfg.IICConfig
	SH1106.Height = cfg.Height
	SH1106.Width = cfg.Width
	SH1106.VccState = cfg.VccState
	err := SH1106.cfg.ReflectFrom(&SH1106)
	if err != nil {
		return err
	}
	return Save()
}
