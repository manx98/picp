package config

import (
	"github.com/go-ini/ini"
	"picp/logger"
)

var wifi = WifiConfig{
	Name: "PICP_202508161714",
}

type WifiConfig struct {
	cfg        *ini.Section `ini:"-"`
	Name       string       `json:"name" ini:"name,omitempty"`
	Enable     bool         `json:"enable" ini:"enable"`
	Pin        int          `json:"pin" ini:"pin" validator:"gte=0,lt=255"`
	SSID       string       `json:"ssid" ini:"ssid" validator:"max=32"`
	Password   string       `json:"password" ini:"password" validator:"min=8"`
	DeviceName string       `json:"device_name" ini:"device_name"`
}

func initWifi() {
	var ok bool
	wifi.cfg, ok = Get("wifi")
	if ok {
		err := StrictMapTo(wifi.cfg, &wifi)
		if err != nil {
			logger.Fatalf("ap config error: %s", err)
		}
	}
}

func GetWifiConfig() WifiConfig {
	cfgLock.RLock()
	defer cfgLock.RUnlock()
	return wifi
}

func SetWifiConfig(cfg *WifiConfig) (err error) {
	err = Validate(cfg)
	if err != nil {
		return
	}
	cfgLock.Lock()
	defer cfgLock.Unlock()
	wifi.Enable = cfg.Enable
	wifi.Pin = cfg.Pin
	wifi.SSID = cfg.SSID
	wifi.Password = cfg.Password
	wifi.DeviceName = cfg.DeviceName
	err = wifi.cfg.ReflectFrom(&wifi)
	if err == nil {
		err = SaveCfg()
	}
	return
}
