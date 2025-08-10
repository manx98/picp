package logger

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

type myWriteSyncer struct {
	ws atomic.Pointer[zapcore.WriteSyncer]
	l  sync.RWMutex
}

func (m *myWriteSyncer) getWs() zapcore.WriteSyncer {
	val := m.ws.Load()
	if val == nil {
		return nil
	}
	return *val
}

func (m *myWriteSyncer) SetWs(ws zapcore.WriteSyncer) {
	if ws == nil || reflect.ValueOf(ws).IsNil() {
		ws = nil
	}
	m.l.Lock()
	val := m.ws.Load()
	m.ws.Store(&ws)
	m.l.Unlock()
	if val != nil {
		err := (*val).Sync()
		zap.L().Warn("failed to sync old log writer", zap.Error(err))
	}
}

func (m *myWriteSyncer) Write(p []byte) (n int, err error) {
	m.l.RLock()
	defer m.l.RUnlock()
	ws := m.getWs()
	if ws == nil {
		return os.Stdout.Write(p)
	} else {
		return ws.Write(p)
	}
}

func (m *myWriteSyncer) Sync() error {
	m.l.Lock()
	defer m.l.Unlock()
	ws := m.getWs()
	if ws == nil {
		return os.Stdout.Sync()
	} else {
		return ws.Sync()
	}
}

var localWs *myWriteSyncer
var logLevel zap.AtomicLevel
var logger *zap.Logger
var logTimeLocation atomic.Pointer[time.Location]
var disableTime atomic.Bool

func init() {
	localWs = &myWriteSyncer{}
	logLevel = zap.NewAtomicLevel()
	logLevel.SetLevel(zapcore.DebugLevel)
	config := zap.NewDevelopmentEncoderConfig()
	config.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		if disableTime.Load() {
			return
		}
		loc := logTimeLocation.Load()
		if loc == nil {
			encoder.AppendString(t.Format("2006-01-02 15:04:05.000"))
		} else {
			encoder.AppendString(t.In(loc).Format("2006-01-02 15:04:05.000"))
		}
	}
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(config), zapcore.Lock(localWs), logLevel)
	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

func GetWitter() io.Writer {
	return localWs
}

func SetLogLevel(level string) (err error) {
	l := zapcore.DebugLevel
	err = l.UnmarshalText([]byte(level))
	if err == nil {
		logLevel.SetLevel(l)
	}
	return
}

func SetLogWriteSyncer(ws zapcore.WriteSyncer) {
	localWs.SetWs(ws)
}

type LumberjackWarp struct {
	writer *lumberjack.Logger
}

func (l *LumberjackWarp) Write(p []byte) (n int, err error) {
	return l.writer.Write(p)
}

func (l *LumberjackWarp) Sync() error {
	return l.writer.Close()
}

func SetLogToFile(logFile string) {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10,
		MaxBackups: 10,
	}
	SetLogWriteSyncer(&LumberjackWarp{writer: lumberJackLogger})
}

func SetTimeLocation(loc string) error {
	location, err := time.LoadLocation(loc)
	if err != nil {
		return errors.WithMessage(err, "load time location")
	}
	logTimeLocation.Store(location)
	return nil
}

func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Debugf(format string, args ...interface{}) {
	logger.Sugar().Debugf(format, args...)
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Infof(format string, args ...interface{}) {
	logger.Sugar().Infof(format, args...)
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func Warnf(format string, args ...interface{}) {
	logger.Sugar().Warnf(format, args...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func Errorf(format string, args ...interface{}) {
	logger.Sugar().Errorf(format, args...)
}

func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Sugar().Fatalf(format, args...)
}

func DisableTime(disable bool) {
	disableTime.Store(disable)
}
