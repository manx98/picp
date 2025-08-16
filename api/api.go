package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"picp/config"
	"picp/driver"
	"picp/logger"
	"picp/utils"
	"strings"
)

func initApi(group *gin.RouterGroup) {
	group.GET("/devices", getDeviceList)
	group.GET("/wifi", getWifiList)
	group.POST("/wifi", connectWifi)
	group.DELETE("/wifi", deleteConnection)
	group.GET("/wifi/config", getWifiConfig)
	group.POST("/wifi/config", setWifiConfig)
	group.GET("/fan", getFanConfig)
	group.POST("/fan", setFanConfig)
	group.GET("/display", getDisplayCfg)
	group.POST("/display", setDisplayCfg)
}

func replayError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": 500,
		"msg":  err.Error(),
	})
}

func replaySuccess(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
	})
}

func getDeviceList(ctx *gin.Context) {
	devices, err := utils.GetDevices()
	if err != nil {
		logger.Error("get devices error", zap.Error(err))
		replayError(ctx, err)
		return
	}
	replaySuccess(ctx, devices)
}

type WifiAPInfo struct {
	utils.WifiAPInfo
	ConnectionUUID string `json:"connection_uuid,omitempty"`
}

func getWifiList(ctx *gin.Context) {
	value := ctx.Query("device")
	if value == "" {
		replayError(ctx, fmt.Errorf("device is empty"))
		return
	}
	wifiList, err := utils.ListWifiAP(value)
	var connections []utils.ConnectionInfo
	if err == nil {
		connections, err = utils.GetConnections()
	}
	ret := make([]*WifiAPInfo, 0, len(wifiList))
	for _, wifi := range wifiList {
		ret = append(ret, &WifiAPInfo{
			WifiAPInfo: wifi,
		})
	}
	if err == nil {
		for _, connection := range connections {
			var info *utils.WifiAPSettingInfo
			if connection.Type == "802-11-wireless" {
				info, err = utils.GetWifiSettingInfo(connection.UUID)
				if err != nil {
					break
				}
				for _, wifi := range ret {
					if info.Is(&wifi.WifiAPInfo) {
						wifi.ConnectionUUID = connection.UUID
					}
				}
			}
		}
	}
	if err != nil {
		logger.Error("get wifi list error", zap.Error(err))
		replayError(ctx, err)
		return
	}
	replaySuccess(ctx, ret)
}

type WifiAPSetting struct {
	SSID           string `json:"SSID"`
	BSSID          string `json:"BSSID,omitempty" validate:"omitempty,mac"`
	Hidden         bool   `json:"hidden"`
	DeviceName     string `json:"device" validate:"required"`
	Password       string `json:"password,omitempty" validate:"omitempty,min=8,max=32"`
	WepKeyType     string `json:"wep_key_type,omitempty" validate:"omitempty,oneof=0 1 2"`
	ConnectionUUID string `json:"connection_uuid,omitempty" validate:"omitempty,uuid"`
	Active         bool   `json:"active"`
}

func connectWifi(ctx *gin.Context) {
	var info WifiAPSetting
	err := ctx.BindJSON(&info)
	if err != nil {
		replayError(ctx, err)
		return
	}
	if info.ConnectionUUID != "" {
		if info.Active {
			err = utils.UpConnection(info.ConnectionUUID)
		} else {
			err = utils.DownConnection(info.ConnectionUUID)
		}
		if err != nil {
			replayError(ctx, err)
			return
		}
	} else {
		if info.SSID != "" && info.BSSID != "" {
			replayError(ctx, fmt.Errorf("SSID and BSSID cannot be set at the same time"))
			return
		}
		err = config.Validate(&info)
		if err != nil {
			replayError(ctx, err)
			return
		}
		var params = utils.WifiAPSetting{
			SSID:       info.SSID,
			BSSID:      info.BSSID,
			Hidden:     info.Hidden,
			DeviceName: info.DeviceName,
			Password:   info.Password,
			WepKeyType: info.WepKeyType,
		}
		if info.SSID != "" {
			params.Name = "WIFI_" + info.SSID
		} else {
			params.Name = "WIFI_" + strings.Replace(info.BSSID, ":", "_", -1)
		}
		err = utils.ConnectWifi(params)
		if err != nil {
			replayError(ctx, err)
			return
		}
	}
	replaySuccess(ctx, nil)
}

func deleteConnection(ctx *gin.Context) {
	connectionUuid := ctx.Query("connection_uuid")
	if connectionUuid == "" {
		replayError(ctx, errors.New("connection_uuid is empty"))
		return
	}
	err := utils.RemoveConnection(connectionUuid)
	if err != nil {
		replayError(ctx, err)
	} else {
		replaySuccess(ctx, nil)
	}
}

func getWifiConfig(ctx *gin.Context) {
	replaySuccess(ctx, config.GetWifiConfig())
}

func setWifiConfig(ctx *gin.Context) {
	var wifiCfg config.WifiConfig
	if err := ctx.ShouldBindJSON(&wifiCfg); err != nil {
		replayError(ctx, err)
		return
	}
	err := driver.SetWifiConfig(&wifiCfg)
	if err != nil {
		replayError(ctx, err)
	} else {
		replaySuccess(ctx, nil)
	}
}

func getFanConfig(ctx *gin.Context) {
	replaySuccess(ctx, config.GetFanCfg())
}

func setFanConfig(ctx *gin.Context) {
	var fanCfg config.FanChanelCfg
	if err := ctx.ShouldBindJSON(&fanCfg); err != nil {
		replayError(ctx, err)
		return
	}
	err := driver.SetFanConfig(&fanCfg)
	if err != nil {
		replayError(ctx, err)
	} else {
		replaySuccess(ctx, nil)
	}
}

func getDisplayCfg(ctx *gin.Context) {
	replaySuccess(ctx, config.GetSH1106Cfg())
}

func setDisplayCfg(ctx *gin.Context) {
	var displayCfg config.SH1106Config
	if err := ctx.ShouldBindJSON(&displayCfg); err != nil {
		replayError(ctx, err)
		return
	}
	err := driver.SetDisplayConfig(&displayCfg)
	if err != nil {
		replayError(ctx, err)
	} else {
		replaySuccess(ctx, nil)
	}
}
