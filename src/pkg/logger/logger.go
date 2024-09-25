package logger

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"project-wraith/src/consts"
	"project-wraith/src/pkg/tools"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Initialize() error
	Warn(message string, args ...interface{})
	Info(message string, args ...interface{})
	Error(message string, args ...interface{})
}

var _ Logger = (*logger)(nil)

type logger struct {
	loggers     map[zapcore.Level]*zap.SugaredLogger
	projectPath string
}

// NewLogger is a function constructor for Logger
func NewLogger(projectPath string) Logger {
	return &logger{
		loggers:     make(map[zapcore.Level]*zap.SugaredLogger),
		projectPath: projectPath,
	}
}

func (l *logger) Initialize() error {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.MessageKey = "message"
	encoderConfig.LevelKey = "severity"

	for _, level := range []zapcore.Level{zap.WarnLevel, zap.InfoLevel, zap.ErrorLevel} {
		logFilePath, err := l.getLogFilePath(level.String())
		if err != nil {
			return err
		}

		cfg := zap.Config{
			Level:             zap.NewAtomicLevelAt(level),
			Development:       false,
			DisableCaller:     true,
			DisableStacktrace: true,
			Encoding:          "json",
			EncoderConfig:     encoderConfig,
			OutputPaths:       []string{"stdout", logFilePath},
			ErrorOutputPaths:  []string{"stderr"},
		}

		zl, err := cfg.Build()
		if err != nil {
			return err
		}

		l.loggers[level] = zl.Sugar()
	}

	if len(l.loggers) < 3 {
		return errors.New("no loggers configured")
	}

	return nil
}

func (l *logger) Warn(message string, args ...interface{}) {
	formattedMessage := fmt.Sprintf(strings.ToLower(message), args...)
	callerInfo := tools.ExtractCallerInfo(consts.LoggerCallerLevel)
	l.loggers[zapcore.WarnLevel].Warnw(formattedMessage, "caller", callerInfo)
}

func (l *logger) Info(message string, args ...interface{}) {
	formattedMessage := fmt.Sprintf(strings.ToLower(message), args...)
	callerInfo := tools.ExtractCallerInfo(consts.LoggerCallerLevel)
	l.loggers[zapcore.InfoLevel].Infow(formattedMessage, "caller", callerInfo)
}

func (l *logger) Error(message string, args ...interface{}) {
	formattedMessage := fmt.Sprintf(strings.ToLower(message), args...)
	callerInfo := tools.ExtractCallerInfo(consts.LoggerCallerLevel)
	l.loggers[zapcore.ErrorLevel].Errorw(formattedMessage, "caller", callerInfo)
}

func (l *logger) getLogFilePath(level string) (string, error) {
	if err := os.MkdirAll(l.projectPath, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create logs directory %e", err)
	}

	return filepath.Join(l.projectPath, fmt.Sprintf("%s.log", level)), nil
}
