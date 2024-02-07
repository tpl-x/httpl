package logger

import (
	"github.com/tpl-x/httpl/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log/slog"
	"os"
)

func NewSlogLogger(config *config.AppConfig) *slog.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	lumberjackLogger := &lumberjack.Logger{
		Filename:   config.Log.FileName,
		MaxSize:    config.Log.MaxSize,
		MaxBackups: config.Log.MaxBackups,
		MaxAge:     config.Log.MaxKeepDays,
		Compress:   config.Log.Compress,
	}
	writeSyncer := zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberjackLogger))
	encoder := zapcore.NewConsoleEncoder(encoderCfg)

	opt := zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewTee(
			zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel),
		)
	})
	zLogger, _ := zap.NewProduction(opt)
	logger := slog.New(zapslog.NewHandler(zLogger.Core(), nil))
	return logger
}
