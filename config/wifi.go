package config

import (
	"errors"
	"github.com/go-ini/ini"
	"picp/logger"
)

var wifiCfg = Wifi{
	Pin:  5,
	Name: "PICP_202508161714",
}

type Wifi struct {
	cfg        *ini.Section `ini:"-"`
	Name       string       `json:"name" ini:"name,omitempty"`
	Enable     bool         `json:"enable" ini:"enable"`
	Pin        int          `json:"pin" ini:"pin" validate:"gte=0,lt=255"`
	SSID       string       `json:"ssid" ini:"ssid" validate:"required"`
	Password   string       `json:"password" ini:"password" validate:"omitempty,min=8,max=32"`
	DeviceName string       `json:"device_name" ini:"device_name"`
}

func (c *Wifi) NeedValidate() bool {
	return c.Enable
}

func initWifi() {
	var ok bool
	wifiCfg.cfg, ok = Get("wifi")
	if ok {
		err := StrictMapTo(wifiCfg.cfg, &wifiCfg)
		if err == nil {

		}
		if err != nil {
			logger.Fatalf("ap config error: %s", err)
		}
	}
}

func GetWifiConfig() Wifi {
	cfgLock.RLock()
	defer cfgLock.RUnlock()
	return wifiCfg
}

func SetWifiConfig(cfg *Wifi) (err error) {
	err = Validate(cfg)
	if err != nil {
		return
	}
	if cfg.SSID == "" && cfg.Enable {
		return errors.New("ssid can not be empty")
	}
	cfgLock.Lock()
	defer cfgLock.Unlock()
	old := wifiCfg
	defer func() {
		if err != nil {
			wifiCfg = old
		}
	}()
	wifiCfg.Enable = cfg.Enable
	wifiCfg.Pin = cfg.Pin
	wifiCfg.SSID = cfg.SSID
	wifiCfg.Password = cfg.Password
	wifiCfg.DeviceName = cfg.DeviceName
	err = wifiCfg.cfg.ReflectFrom(&wifiCfg)
	if err == nil {
		err = SaveCfg()
	}
	return
}
