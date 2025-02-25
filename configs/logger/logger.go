package logger

import (
	"go-rest-api-boilerplate/types"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

type LoggerZap struct {
	*zap.Logger
}

// NewLoggerZap creates a new LoggerZap
func NewLogger(config types.LoggerSetting) *LoggerZap {
	logLevel := config.LogLevel
	//debug -> info -> warn -> error -> fatal -> panic

	var level zapcore.Level
	switch logLevel {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	encoder := getEndcoderLog()
	hook := lumberjack.Logger{
		Filename:   config.FileLogName,
		MaxSize:    config.MaxSize,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
		MaxBackups: config.MaxBackups,
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)),
		level,
	)

	//logger := zap.New(core, zap.AddCaller())
	return &LoggerZap{zap.New(core, zap.AddCaller())}

}

func getEndcoderLog() zapcore.Encoder {
	encodeConfig := zap.NewProductionEncoderConfig()

	//
	encodeConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	//ts -> Time
	encodeConfig.TimeKey = "time"

	//from info INFO
	encodeConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	//"caller" : "cli/main.go:18"
	encodeConfig.EncodeCaller = zapcore.ShortCallerEncoder

	return zapcore.NewJSONEncoder(encodeConfig)
}
