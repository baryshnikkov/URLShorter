package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func New(pathFile string) *zap.Logger {
	file, err := os.OpenFile(pathFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic("cannot open log file")
	}

	fileSync := zapcore.AddSync(file)
	consoleSync := zapcore.AddSync(os.Stdout)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	encoder := zapcore.NewJSONEncoder(encoderConfig)

	logLevel := zapcore.InfoLevel

	core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(fileSync, consoleSync), logLevel)

	logger := zap.New(core)
	defer logger.Sync()

	return logger
}
