package config

import (
	"errors"
	"flag"
	"fmt"
	"github.com/go-ini/ini"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"os"
	"picp/logger"
	"reflect"
	"strings"
)

var vid *validator.Validate

var cfgPath = flag.String("c", "picp.ini", "config file path")
var rootCfg *ini.File

type NeedValidateFunc interface {
	NeedValidate() bool
}

func Init() {
	vid = validator.New()
	vid.RegisterTagNameFunc(func(field reflect.StructField) string {
		if jsonTag, ok := field.Tag.Lookup("ini"); ok {
			return strings.SplitN(jsonTag, ",", 2)[0]
		} else {
			return field.Name
		}
	})
	var err error
	rootCfg, err = ini.Load(*cfgPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			rootCfg = ini.Empty()
		} else {
			logger.Fatal("load config failed", zap.Error(err))
		}
	}
	initLogger()
	initCommon()
	initSH1106()
	initFan()
	initWifi()
}

func Save() error {
	err := rootCfg.SaveTo(*cfgPath)
	if err != nil {
		return fmt.Errorf("save config failed: %w", err)
	}
	return nil
}

func Get(section string) (*ini.Section, bool) {
	return rootCfg.Section(section), rootCfg.HasSection(section)
}

func StrictMapTo(cfg *ini.Section, obj interface{}) error {
	err := cfg.StrictMapTo(obj)
	if err != nil {
		return fmt.Errorf("invalid config %s: %w", cfg.Name(), err)
	}
	err = Validate(obj)
	if err != nil {
		return fmt.Errorf("invalid config %s: %w", cfg.Name(), err)
	}
	logger.Debug("map config success", zap.String("name", cfg.Name()), zap.Any("cfg", obj))
	return nil
}
func Validate(config interface{}) (err error) {
	if obj, ok := config.(NeedValidateFunc); !ok || obj.NeedValidate() {
		err = vid.Struct(config)
	}
	return err
}

func initLogger() {
	section, _ := Get("common")
	logLevel := section.Key("log_level").MustString("info")
	logFile := section.Key("log_file").MustString("")
	err := logger.SetLogLevel(logLevel)
	if err != nil {
		logger.Fatal("set log level failed", zap.String("log level", logLevel), zap.Error(err))
	}
	if logFile != "" {
		logger.SetLogToFile(logFile)
	}
}
