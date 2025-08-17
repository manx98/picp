package driver

import (
	"context"
	"github.com/stianeikeland/go-rpio/v4"
	"picp/config"
	"picp/utils"
	"strings"
	"time"
)

var wifiRunner *utils.Runner

func wifiInit(ctx context.Context) {
	wifiRunner = utils.NewRunner(ctx, wifiHandler)
	wifiRunner.Start()
}

type WifiInvoker struct {
	cfg           *config.Wifi
	enabled       bool
	notifyTimeout time.Time
	ctx           context.Context
}

func (i *WifiInvoker) startAp() {
	i.showNotify("Connect...")
	err := startWifiAp(i.cfg)
	if err == nil {
		i.enabled = true
		i.showNotify("Connect success")
	} else {
		i.showNotify(strings.Split(err.Error(), "\n")...)
	}
}

func (i *WifiInvoker) stopAp() {
	i.showNotify("Stopping...")
	err := stopWifiAp(i.cfg)
	if err == nil {
		i.enabled = false
		i.showNotify("Stop success")
	} else {
		i.showNotify(strings.Split(err.Error(), "\n")...)
	}
}

func (i *WifiInvoker) checkNotify(force bool) bool {
	if !i.notifyTimeout.IsZero() {
		if force || time.Now().After(i.notifyTimeout) {
			i.notifyTimeout = time.Time{}
			StatusShowEnable(true)
		}
		return true
	}
	return false
}

func (i *WifiInvoker) showNotify(msg ...string) {
	StatusShowEnable(false)
	i.notifyTimeout = time.Now().Add(time.Second * 3)
	DisplayVerticalAlign(msg...)
}

func (i *WifiInvoker) toggleAp() {
	if i.enabled {
		i.stopAp()
	} else {
		i.startAp()
	}
}

func (i *WifiInvoker) Run() {
	_ = stopWifiAp(i.cfg)
	if !i.cfg.Enable {
		return
	}
	wifiPin := rpio.Pin(i.cfg.Pin)
	wifiPin.Input()
	timer := time.NewTicker(time.Millisecond * 100)
	defer func() {
		timer.Stop()
		_ = stopWifiAp(i.cfg)
		StatusShowEnable(true)
	}()
	lastApPress := false
	for {
		select {
		case <-timer.C:
			press := wifiPin.Read() == rpio.High
			if press != lastApPress {
				if lastApPress {
					if !i.checkNotify(true) {
						i.toggleAp()
					}
				} else {
					i.checkNotify(false)
				}
				lastApPress = !lastApPress
			} else {
				i.checkNotify(false)
			}
		case <-i.ctx.Done():
			return
		}
	}
}

func NewWifiInvoker(ctx context.Context, cfg *config.Wifi) *WifiInvoker {
	invoker := &WifiInvoker{cfg: cfg, ctx: ctx}
	return invoker
}

func wifiHandler(ctx context.Context) {
	cfg := config.GetWifiConfig()
	_ = stopWifiAp(&cfg)
	if !cfg.Enable {
		return
	}
	NewWifiInvoker(ctx, &cfg).Run()
}

func startWifiAp(cfg *config.Wifi) error {
	err := utils.ForceScan()
	if err != nil {
		return err
	}
	err = utils.ConnectWifi(utils.WifiAPSetting{
		Name:       cfg.Name,
		SSID:       cfg.SSID,
		DeviceName: cfg.DeviceName,
		Password:   cfg.Password,
	})
	if err != nil {
		_ = utils.RemoveConnectionByID(cfg.Name)
	}
	return err
}

func stopWifiAp(cfg *config.Wifi) error {
	return utils.RemoveConnectionByID(cfg.Name)
}

func SetWifiConfig(wifi *config.Wifi) error {
	err := config.SetWifiConfig(wifi)
	if err != nil {
		return err
	}
	_ = wifiRunner.Stop(context.Background())
	wifiRunner.Start()
	return nil
}

func closeWifi() {
	_ = wifiRunner.Stop(context.Background())
}
