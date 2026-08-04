package logger

import (
	"boilerplate/config"
	"fmt"
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Sugared struct {
	*zap.SugaredLogger
}

var Zap *Sugared

func New() *Sugared {
	conf, _ := config.NewConfig()

	level := zap.DebugLevel
	if conf != nil && conf.App.Log.Level != "" {
		levelFromEnv, err := zapcore.ParseLevel(conf.App.Log.Level)
		if err != nil {
			log.Println(fmt.Errorf("invalid level, defaulting to DEBUG: %w", err))
		} else {
			level = levelFromEnv
		}
	}

	var encoderConfig zapcore.EncoderConfig
	if level == zap.DebugLevel {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		encoderConfig = zap.NewProductionEncoderConfig()
	}

	encoderConfig.MessageKey = "message"
	encoderConfig.LevelKey = "level"
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.CallerKey = "caller"
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	zapConfig := zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Encoding:          "json",
		EncoderConfig:     encoderConfig,
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
	}

	Zap = &Sugared{zap.Must(zapConfig.Build()).Sugar()}
	return Zap
}

func (l *Sugared) Printf(format string, args ...interface{}) {
	l.Infof(format, args...)
}
