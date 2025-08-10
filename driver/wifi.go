package driver

import (
	"github.com/stianeikeland/go-rpio/v4"
	"picp/config"
)

var wifiPin rpio.Pin

func wifiInit() {
	if config.WIFI.Enable {
		wifiPin = rpio.Pin(config.WIFI.Pin)
		wifiPin.Input()
	}
}

func isWifi() bool {
	if config.WIFI.Enable {
		return wifiPin.Read() == rpio.High
	}
	return false
}
