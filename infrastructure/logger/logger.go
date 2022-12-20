package logger

import (
	"errors"

	"go.uber.org/zap"
)

var ErrLoadingLogger = errors.New("ロガーの生成中にエラーが発生しました")

type Logger struct {
	logger *zap.SugaredLogger
}

func (l *Logger) Info(args ...any) {
	l.logger.Info(args)
}

func (l *Logger) Debug(args ...any) {
	l.logger.Debug(args)
}

func (l *Logger) Warn(args ...any) {
	l.logger.Warn(args)
}

func (l *Logger) Error(args ...any) {
	l.logger.Error(args)
}

func (l *Logger) Fatal(args ...any) {
	l.logger.Fatal(args)
}

func NewLogger(isDevelopment bool) (*Logger, error) {
	var err error

	var loggerBase *zap.Logger
	if isDevelopment {
		loggerBase, err = zap.NewDevelopment()
	} else {
		loggerBase, err = zap.NewProduction()
	}
	if err != nil {
		return nil, ErrLoadingLogger
	}
	defer loggerBase.Sync()

	return &Logger{
		logger: loggerBase.Sugar(),
	}, nil
}
