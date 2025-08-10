package utils

import (
	"go.uber.org/zap"
	"picp/logger"
	"testing"
)

func TestGetConnections(t *testing.T) {
	connections, err := GetConnections()
	if err != nil {
		t.Fatal("GetConnections", err)
	}
	for _, info := range connections {
		if info.Type == "802-11-wireless" {
			wifiSetting, err := GetWifiSettingInfo(info.UUID)
			if err != nil {
				t.Fatal("GetWifiSettingInfo", err)
			}
			logger.Info("GetWifiSettingInfo", zap.Any("wifiSetting", wifiSetting))
		}
		logger.Info("GetConnections", zap.Any("info", info))
	}
	devices, err := GetDevices()
	if err != nil {
		t.Fatal("GetDevices", err)
	}
	for _, info := range devices {
		logger.Info("GetDevices", zap.Any("info", info))
	}
	wifiConnections, err := ListWifiAP("wlan0")
	if err != nil {
		t.Fatal("ListWifiAP", err)
	}
	for _, info := range wifiConnections {
		logger.Info("ListWifiAP", zap.Any("info", info))
	}
	info, err := GetRadioInfo()
	if err != nil {
		t.Fatal("GetRadioInfo", err)
	}
	logger.Info("GetRadioInfo", zap.Any("info", info))
}
