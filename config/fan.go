package config

import (
	"github.com/go-ini/ini"
	"picp/logger"
)

var fan = FanChanelCfg{
	Enable:  false,
	Pin:     12,
	Speed:   60,
	MaxTemp: 50,
	MinTemp: 45,
}

type FanChanelCfg struct {
	cfg     *ini.Section `ini:"-"`
	Enable  bool         `json:"enable" ini:"enable,omitempty"`
	Pin     int          `json:"pin" ini:"pin,omitempty" validate:"oneof=12 13 40 41 45 18 19"`
	Speed   int          `json:"speed" ini:"speed,omitempty" validate:"gt=0,lte=100"`
	MinTemp float32      `json:"min_temp" ini:"min_temp,omitempty" validate:"gte=0,ltfield=MaxTemp"`
	MaxTemp float32      `json:"max_temp" ini:"max_temp,omitempty" validate:"gte=0"`
}

func (c *FanChanelCfg) NeedValidate() bool {
	return c.Enable
}

func initFan() {
	var ok bool
	fan.cfg, ok = Get("fan")
	if ok {
		if err := StrictMapTo(fan.cfg, &fan); err != nil {
			logger.Fatalf("fan config error: %s", err)
		}
	}
}

func GetFanCfg() FanChanelCfg {
	cfgLock.RLock()
	defer cfgLock.RUnlock()
	return fan
}

func SetFanCfg(cfg *FanChanelCfg) (err error) {
	err = Validate(cfg)
	if err != nil {
		return
	}
	cfgLock.Lock()
	defer cfgLock.Unlock()
	old := fan
	defer func() {
		if err != nil {
			fan = old
		}
	}()
	fan.Enable = cfg.Enable
	fan.Pin = cfg.Pin
	fan.Speed = cfg.Speed
	fan.MinTemp = cfg.MinTemp
	fan.MaxTemp = cfg.MaxTemp
	err = fan.cfg.ReflectFrom(&fan)
	if err == nil {
		return SaveCfg()
	}
	return
}
