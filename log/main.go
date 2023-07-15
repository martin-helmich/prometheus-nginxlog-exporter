package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	zap *zap.SugaredLogger
}

func New(logLevel, logFormat string) (*Logger, error) {
	level, err := zap.ParseAtomicLevel(logLevel)
	if err != nil {
		return nil, err
	}

	config := zap.NewProductionConfig()
	config.Level = level
	config.Encoding = logFormat
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// build logger
	log, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}

	return &Logger{
		zap: log.Sugar(),
	}, nil
}

func (log *Logger) Print(args ...interface{}) {
	log.zap.Info(args...)
}

func (log *Logger) Debug(msg string, fields ...zap.Field) {
	log.zap.Debug(msg, fields)
}

func (log *Logger) Info(msg string, fields ...zap.Field) {
	log.zap.Info(msg, fields)
}

func (log *Logger) Warn(msg string, fields ...zap.Field) {
	log.zap.Warn(msg, fields)
}

func (log *Logger) Error(msg string, fields ...zap.Field) {
	log.zap.Error(msg, fields)
}

func (log *Logger) Fatal(args ...interface{}) {
	log.zap.Fatal(args...)
}

func (log *Logger) Panic(args ...interface{}) {
	log.zap.Panic(args...)
}

func (log *Logger) Printf(template string, args ...interface{}) {
	log.zap.Infof(template, args...)
}

func (log *Logger) Debugf(template string, args ...interface{}) {
	log.zap.Debugf(template, args...)
}
func (log *Logger) Infof(template string, args ...interface{}) {
	log.zap.Infof(template, args...)
}

func (log *Logger) Warnf(template string, args ...interface{}) {
	log.zap.Warnf(template, args...)
}

func (log *Logger) Errorf(template string, args ...interface{}) {
	log.zap.Errorf(template, args...)
}

func (log *Logger) Fatalf(template string, args ...interface{}) {
	log.zap.Fatalf(template, args...)
}

func (log *Logger) Panicf(template string, args ...interface{}) {
	log.zap.Panicf(template, args...)
}

func (log *Logger) Println(args ...interface{}) {
	log.zap.Infoln(args...)
}

func (log *Logger) Debugln(args ...interface{}) {
	log.zap.Debugln(args...)
}

func (log *Logger) Infoln(args ...interface{}) {
	log.zap.Infoln(args...)
}

func (log *Logger) Warnln(args ...interface{}) {
	log.zap.Warnln(args...)
}

func (log *Logger) Errorln(args ...interface{}) {
	log.zap.Errorln(args...)
}

func (log *Logger) Fatalln(args ...interface{}) {
	log.zap.Fatalln(args...)
}

func (log *Logger) Panicln(args ...interface{}) {
	log.zap.Panicln(args...)
}
