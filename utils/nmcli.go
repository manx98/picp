package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net"
	"os/exec"
	"picp/logger"
	"strconv"
	"strings"
)

type NMActiveConnectionState int

func (s *NMActiveConnectionState) String() string {
	index := int(*s)
	if index >= len(nmActiveConnectionStateString) {
		return fmt.Sprintf("unknown(%d)", *s)
	}
	return nmActiveConnectionStateString[index]
}

func (s *NMActiveConnectionState) FromString(value string) bool {
	for i, name := range nmActiveConnectionStateString {
		if name == value {
			*s = NMActiveConnectionState(i)
			return true
		}
	}
	*s = NMActiveConnectionStateUnknown
	return false
}

var nmActiveConnectionStateString = []string{"unknown", "activating", "activated", "deactivating", "deactivated"}

const (
	NMActiveConnectionStateUnknown = NMActiveConnectionState(iota)
	NMActiveConnectionStateActivating
	NMActiveConnectionStateActivated
	NMActiveConnectionStateDeactivating
	NMActiveConnectionStateDeactivated
)

type ConnectionInfo struct {
	Name                string                  `json:"name,omitempty"`
	UUID                string                  `json:"uuid,omitempty"`
	Type                string                  `json:"type,omitempty"`
	Device              string                  `json:"device,omitempty"`
	State               NMActiveConnectionState `json:"state,omitempty"`
	AutoConnect         bool                    `json:"auto_connect"`
	AutoConnectPriority int                     `json:"auto_connect_priority,omitempty"`
}

func parseLineValues(line string) ([]string, error) {
	var strBuf strings.Builder
	var ret []string
	lineRunes := []rune(line)
	directWrite := false
	for i, c := range lineRunes {
		if directWrite {
			strBuf.WriteRune(c)
			directWrite = false
		} else if c == ':' {
			ret = append(ret, strBuf.String())
			strBuf.Reset()
		} else if c == '\\' {
			if i+1 < len(line) {
				if lineRunes[i+1] == '\\' || lineRunes[i+1] == ':' {
					directWrite = true
				} else {
					return nil, fmt.Errorf("invalid escape sequence at %d", i)
				}
			} else {
				return nil, fmt.Errorf("invalid escape sequence at %d", i)
			}
		} else {
			strBuf.WriteRune(c)
		}
	}
	ret = append(ret, strBuf.String())
	return ret, nil
}

func runCmd(name string, arg ...string) (string, error) {
	command := exec.Command(name, arg...)
	output, err := command.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ERROR %s: %w", output, err)
	}
	return string(output), nil
}

func GetConnections() ([]ConnectionInfo, error) {
	output, err := runCmd("nmcli", "-t", "-f", "name,uuid,type,device,state,autoconnect,autoconnect-priority", "connection", "show")
	if err != nil {
		return nil, err
	}
	logger.Debug("nmcli connections", zap.String("output", output))
	var ret []ConnectionInfo
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		values, err := parseLineValues(line)
		if err != nil {
			logger.Warn("failed to parse nmcli connections", zap.Error(err), zap.String("line", line))
			return nil, fmt.Errorf("parse nmcli connections line : %w", err)
		}
		if len(values) == 7 {
			var obj ConnectionInfo
			obj.Name = values[0]
			obj.UUID = values[1]
			obj.Type = values[2]
			obj.Device = values[3]
			obj.State.FromString(values[4])
			obj.AutoConnect = values[5] == "yes"
			obj.AutoConnectPriority, _ = strconv.Atoi(values[6])
			ret = append(ret, obj)
		} else {
			logger.Warn("nmcli connection info length mismatch", zap.String("line", line), zap.Strings("values", values))
			return nil, errors.New("invalid connection info line")
		}
	}
	return ret, nil
}

const (
	DeviceKey       = "GENERAL.DEVICE:"
	TypeKey         = "GENERAL.TYPE:"
	StateKey        = "GENERAL.STATE:"
	HardwareAddrKey = "GENERAL.HWADDR:"
	IPV4AddrKey     = "IP4.ADDRESS["
	IPV4GatewayKey  = "IP4.GATEWAY:"
	IPV6AddrKey     = "IP6.ADDRESS["
	IPV6GatewayKey  = "IP6.GATEWAY:"
	ConnectUUIDKey  = "GENERAL.CON-UUID:"
	DeviceUpKey     = "INTERFACE-FLAGS.UP:"
)

type NMDeviceState int

func (s *NMDeviceState) String() string {
	for i, value := range nmDeviceStateValues {
		if value == *s {
			return nmDeviceStateString[i]
		}
	}
	return fmt.Sprintf("unknown(%d)", *s)
}

func (s *NMDeviceState) FromString(value string) bool {
	for i, stateStr := range nmDeviceStateString {
		if stateStr == value {
			*s = nmDeviceStateValues[i]
			return true
		}
	}
	*s = NMDeviceStateUnknown
	return false
}

func (s *NMDeviceState) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

var nmDeviceStateValues = []NMDeviceState{
	NMDeviceStateUnmanaged,
	NMDeviceStateUnavailable,
	NMDeviceStateDisconnected,
	NMDeviceStatePrepare,
	NMDeviceStateConfig,
	NMDeviceStateNeedAuth,
	NMDeviceStateIpConfig,
	NMDeviceStateIpCheck,
	NMDeviceStateActivated,
	NMDeviceStateDeactivating,
	NMDeviceStateFailed,
	NMDeviceStateUnknown,
	NMDeviceStateSecondaries,
}

var nmDeviceStateString = []string{
	"unmanaged",
	"unavailable",
	"disconnected",
	"connecting (prepare)",
	"connecting (configuring)",
	"connecting (need authentication)",
	"connecting (getting IP configuration)",
	"connecting (checking IP connectivity)",
	"connected",
	"deactivating",
	"connection",
	"unknown",
	"connecting (starting secondary connections)",
}

const (
	NMDeviceStateUnknown      NMDeviceState = 0
	NMDeviceStateUnmanaged    NMDeviceState = 10
	NMDeviceStateUnavailable  NMDeviceState = 20
	NMDeviceStateDisconnected NMDeviceState = 30
	NMDeviceStatePrepare      NMDeviceState = 40
	NMDeviceStateConfig       NMDeviceState = 50
	NMDeviceStateNeedAuth     NMDeviceState = 60
	NMDeviceStateIpConfig     NMDeviceState = 70
	NMDeviceStateIpCheck      NMDeviceState = 80
	NMDeviceStateSecondaries  NMDeviceState = 90
	NMDeviceStateActivated    NMDeviceState = 100
	NMDeviceStateDeactivating NMDeviceState = 110
	NMDeviceStateFailed       NMDeviceState = 120
)

type DeviceInfo struct {
	Up             bool          `json:"up"`
	Device         string        `json:"device"`
	ConnectionUUID string        `json:"connection_uuid"`
	Type           string        `json:"type"`
	State          NMDeviceState `json:"state"`
	HardwareAddr   string        `json:"hardware_addr"`
	IPV4           []net.IP      `json:"ipv4"`
	IPV4Gateway    net.IP        `json:"ipv4_gateway"`
	IPV6           []net.IP      `json:"ipv6"`
	IPV6Gateway    net.IP        `json:"ipv6_gateway"`
}

func getLineStringValue(line, key string, value *string) bool {
	if strings.HasPrefix(line, key) {
		*value = strings.TrimSpace(line[len(key):])
		return true
	}
	return false
}

func appendLineIPValue(line, key string, value *[]net.IP) bool {
	if strings.HasPrefix(line, key) {
		index := strings.Index(line, ":")
		if index != -1 {
			ip := strings.TrimSpace(line[index+1:])
			ipAddr, _, err := net.ParseCIDR(ip)
			if err != nil {
				logger.Warn("failed to resolve ip", zap.Error(err), zap.String("ip", ip))
			} else {
				*value = append(*value, ipAddr)
			}
		} else {
			logger.Warn("invalid ip address", zap.String("line", line))
		}
		return true
	}
	return false
}

func getLineIPValue(line, key string, value *net.IP) bool {
	if strings.HasPrefix(line, key) {
		line = line[len(key):]
		ipAddr := net.ParseIP(line)
		*value = ipAddr
		return true
	}
	return false
}
func getLineStateValue(line string, value *NMDeviceState) bool {
	if strings.HasPrefix(line, StateKey) {
		valueStr := line[len(StateKey):]
		stateIndex := strings.Index(valueStr, " ")
		if stateIndex != -1 {
			valueStr = valueStr[:stateIndex]
			state, err := strconv.Atoi(valueStr)
			if err != nil {
				logger.Warn("failed to resolve state", zap.Error(err), zap.String("state", line))
			} else {
				*value = NMDeviceState(state)
				return true
			}
		} else {
			logger.Warn("invalid state", zap.String("state", line))
		}
		*value = NMDeviceStateUnknown
		return true
	}
	return false
}

func getLineBoolValue(line, key string, value *bool) bool {
	if strings.HasPrefix(line, key) {
		*value = line[len(key):] == "yes"
		return true
	}
	return false
}

func GetDevices() ([]DeviceInfo, error) {
	cmd, err := runCmd("nmcli", "-t", "-f", "INTERFACE-FLAGS.UP,GENERAL.DEVICE,GENERAL.TYPE,GENERAL.STATE,GENERAL.HWADDR,IP4.ADDRESS,IP4.GATEWAY,IP6.ADDRESS,IP6.GATEWAY,GENERAL.CON-UUID", "device", "show")
	if err != nil {
		return nil, fmt.Errorf("run nmcli device show failed : %w", err)
	}
	logger.Info("nmcli device show", zap.String("out", cmd))
	lines := strings.Split(cmd, "\n")
	var ret []DeviceInfo
	var info DeviceInfo
	for _, line := range lines {
		if len(line) == 0 {
			ret = append(ret, info)
			info = DeviceInfo{}
		} else {
			if getLineStringValue(line, DeviceKey, &info.Device) {
			} else if getLineStringValue(line, TypeKey, &info.Type) {
			} else if getLineStateValue(line, &info.State) {
			} else if getLineStringValue(line, HardwareAddrKey, &info.HardwareAddr) {
			} else if appendLineIPValue(line, IPV4AddrKey, &info.IPV4) {
			} else if getLineIPValue(line, IPV4GatewayKey, &info.IPV4Gateway) {
			} else if appendLineIPValue(line, IPV6AddrKey, &info.IPV6) {
			} else if getLineIPValue(line, IPV6GatewayKey, &info.IPV6Gateway) {
			} else if getLineStringValue(line, ConnectUUIDKey, &info.ConnectionUUID) {
			} else if getLineBoolValue(line, DeviceUpKey, &info.Up) {
			}
		}
	}
	return ret, nil
}

func RemoveConnection(uuid string) error {
	out, err := runCmd("nmcli", "connection", "delete", "uuid", uuid)
	if err != nil {
		logger.Warn("failed to delete connection", zap.String("UUID", uuid), zap.Error(err))
		return err
	} else {
		logger.Debug("delete connection success", zap.String("UUID", uuid), zap.String("out", out))
	}
	return nil
}

var ErrPasswordsRequired = errors.New("passwords required")

func UpConnection(uuid string) error {
	out, err := runCmd("nmcli", "connection", "up", "uuid", uuid)
	if err != nil {
		logger.Warn("failed to connection", zap.String("uuid", uuid), zap.Error(err))
		if strings.Contains(err.Error(), "Passwords or encryption keys are required to access the wireless network") {
			return ErrPasswordsRequired
		}
		return err
	} else {
		logger.Debug("connection", zap.String("uuid", uuid), zap.String("out", out))
	}
	return nil
}

func DownConnection(uuid string) error {
	out, err := runCmd("nmcli", "connection", "down", "uuid", uuid)
	if err != nil {
		logger.Warn("failed to down connection", zap.String("uuid", uuid), zap.Error(err))
		return err
	} else {
		logger.Debug("disconnect connection", zap.String("uuid", uuid), zap.String("out", out))
	}
	return nil
}

func ConnectDevice(device string) error {
	out, err := runCmd("nmcli", "device", "connect", device)
	if err != nil {
		logger.Warn("failed to connect device", zap.String("device", device), zap.Error(err))
		return err
	} else {
		logger.Debug("connect device", zap.String("device", device), zap.String("out", out))
	}
	return nil
}

func DisconnectDevice(device string) error {
	out, err := runCmd("nmcli", "device", "disconnect", device)
	if err != nil {
		logger.Warn("failed to disconnect device", zap.String("device", device), zap.Error(err))
		return err
	} else {
		logger.Debug("disconnect device", zap.String("device", device), zap.String("out", out))
	}
	return nil
}

type WifiAPInfo struct {
	SSID     string   `json:"SSID,omitempty"`
	BSSID    string   `json:"BSSID,omitempty"`
	Mode     string   `json:"mode,omitempty"`
	Chan     int      `json:"chan,omitempty"`
	Freq     string   `json:"freq,omitempty"`
	Rate     string   `json:"rate,omitempty"`
	Signal   int      `json:"signal,omitempty"`
	Security []string `json:"security,omitempty"`
	Active   bool     `json:"active,omitempty"`
	Device   string   `json:"device"`
}

func ListWifiAP(ifName string) ([]WifiAPInfo, error) {
	cmd, err := runCmd("nmcli", "-t", "-f", "SSID,BSSID,MODE,CHAN,FREQ,RATE,SIGNAL,SECURITY,ACTIVE,DEVICE", "device", "wifi", "list", "ifname", ifName, "--rescan", "yes")
	if err != nil {
		logger.Warn("run nmcli device wifi list failed", zap.String("ifName", ifName), zap.Error(err))
		return nil, err
	}
	logger.Debug("nmcli device wifi list", zap.String("ifName", ifName), zap.String("out", cmd))
	lines := strings.Split(cmd, "\n")
	var ret []WifiAPInfo
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		values, err := parseLineValues(line)
		if err != nil {
			logger.Warn("failed to parse line", zap.String("line", line), zap.Error(err))
			continue
		}
		if len(values) != 10 {
			logger.Warn("invalid line", zap.String("line", line))
			continue
		}
		info := WifiAPInfo{
			SSID:   values[0],
			BSSID:  values[1],
			Mode:   values[2],
			Freq:   values[4],
			Rate:   values[5],
			Active: values[8] == "yes",
			Device: values[9],
		}
		if values[7] != "" {
			info.Security = strings.Split(values[7], " ")
		}
		info.Chan, _ = strconv.Atoi(values[3])
		info.Signal, _ = strconv.Atoi(values[6])
		ret = append(ret, info)
	}
	return ret, nil
}

type WifiAPSetting struct {
	Name       string `json:"-"`
	SSID       string `json:"SSID,omitempty"`
	Hidden     bool   `json:"hidden,omitempty"`
	BSSID      string `json:"BSSID,omitempty"`
	DeviceName string `json:"device_name,omitempty"`
	Password   string `json:"password,omitempty"`
	WepKeyType string `json:"wep_key_type,omitempty"`
}

func (w *WifiAPSetting) toCmd() (cmd []string) {
	cmd = append(cmd, "device", "wifi", "connect", w.SSID)
	if w.Hidden {
		cmd = append(cmd, "hidden", "yes")
	}
	if w.BSSID != "" {
		cmd = append(cmd, "bssid", w.BSSID)
	}
	if w.Password != "" {
		cmd = append(cmd, "password", w.Password)
	}
	if w.DeviceName != "" {
		cmd = append(cmd, "ifname", w.DeviceName)
	}
	if w.WepKeyType != "" {
		cmd = append(cmd, "wep-key-type", w.WepKeyType)
	}
	if w.Name != "" {
		cmd = append(cmd, "name", w.Name)
	}
	return
}

func ConnectWifi(info WifiAPSetting) error {
	out, err := runCmd("nmcli", info.toCmd()...)
	if err != nil {
		logger.Warn("run nmcli device wifi connect failed", zap.Error(err))
		return err
	} else {
		if strings.Contains(out, "Secrets were required, but not provided.") {
			logger.Warn("wifi password required", zap.String("out", out))
			_, _ = runCmd("nmcli", "connection", "delete", info.Name)
			return ErrPasswordsRequired
		}
		logger.Debug("connect wifi success", zap.String("out", out))
	}
	return nil
}

const (
	wifiApSettingSSIDKey       = "802-11-wireless.ssid:"
	wifiApSettingBSSIDKey      = "802-11-wireless.bssid:"
	wifiApSettingHiddenKey     = "802-11-wireless.hidden:"
	wifiApSettingDeviceNameKey = "connection.interface-name:"
	wifiApSettingWepKeyTypeKey = "802-11-wireless-security.wep-key-type:"
	wifiApSettingModeKey       = "802-11-wireless.mode:"
)

type WifiAPSettingInfo struct {
	SSID       string `json:"SSID,omitempty"`
	BSSID      string `json:"BSSID,omitempty"`
	Hidden     bool   `json:"hidden,omitempty"`
	DeviceName string `json:"device_name,omitempty"`
	Password   string `json:"password,omitempty"`
	WepKeyType string `json:"wep_key_type,omitempty"`
	Mode       string `json:"mode,omitempty"`
}

func (w *WifiAPSettingInfo) Is(ap *WifiAPInfo) bool {
	if w.BSSID != "" && ap.BSSID != w.BSSID {
		return false
	}
	if ap.SSID != w.SSID {
		return false
	}
	if w.DeviceName != "" && ap.Device != w.DeviceName {
		return false
	}
	if ap.Mode != "Infra" || w.Mode != "infrastructure" {
		return false
	}
	return true
}

func GetWifiSettingInfo(connectUUID string) (*WifiAPSettingInfo, error) {
	cmd, err := runCmd("nmcli", "-t", "-f", "connection.interface-name,802-11-wireless,802-11-wireless-security.wep-key-type", "connection", "show", connectUUID)
	if err != nil {
		logger.Warn("run nmcli connection show failed", zap.String("out", cmd))
		return nil, err
	}
	logger.Debug("get wifi setting info", zap.String("out", cmd))
	lines := strings.Split(cmd, "\n")
	var setting WifiAPSettingInfo
	for _, line := range lines {
		if getLineStringValue(line, wifiApSettingSSIDKey, &setting.SSID) {
			continue
		} else if getLineStringValue(line, wifiApSettingBSSIDKey, &setting.BSSID) {
			continue
		} else if getLineBoolValue(line, wifiApSettingHiddenKey, &setting.Hidden) {
			continue
		} else if getLineStringValue(line, wifiApSettingDeviceNameKey, &setting.DeviceName) {
			continue
		} else if getLineStringValue(line, wifiApSettingWepKeyTypeKey, &setting.WepKeyType) {
			if setting.WepKeyType == "unknown" {
				setting.WepKeyType = ""
			}
			continue
		} else if getLineStringValue(line, wifiApSettingModeKey, &setting.Mode) {
			continue
		}
	}
	return &setting, nil
}
func RadioSwitch(isWifi, state bool) error {
	args := []string{"radio"}
	if isWifi {
		args = append(args, "wifi")
	} else {
		args = append(args, "wwan")
	}
	if state {
		args = append(args, "on")
	} else {
		args = append(args, "off")
	}
	cmd, err := runCmd("nmcli", args...)
	if err != nil {
		logger.Warn("run nmcli radio wifi failed", zap.String("out", cmd), zap.Strings("args", args))
		return err
	}
	return nil
}

type RadioInfo struct {
	WifiHW bool `json:"wifiHW,omitempty"`
	Wifi   bool `json:"wifi,omitempty"`
	WwanHW bool `json:"wwanHW,omitempty"`
	Wwan   bool `json:"wwan,omitempty"`
}

func GetRadioInfo() (*RadioInfo, error) {
	cmd, err := runCmd("nmcli", "-t", "radio", "all")
	if err != nil {
		logger.Warn("run nmcli radio failed", zap.String("out", cmd))
		return nil, err
	}
	logger.Debug("get radio info", zap.String("out", cmd))
	line := strings.Trim(cmd, "\n")
	values, err := parseLineValues(line)
	if err != nil {
		logger.Warn("failed to parse line", zap.String("line", line), zap.Error(err))
		return nil, err
	}
	if len(values) != 4 {
		logger.Warn("invalid line", zap.String("line", line))
		return nil, errors.New("invalid result")
	}
	return &RadioInfo{
		WifiHW: values[0] == "enabled",
		Wifi:   values[1] == "enabled",
		WwanHW: values[2] == "enabled",
		Wwan:   values[3] == "enabled",
	}, nil
}
