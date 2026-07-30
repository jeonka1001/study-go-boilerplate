package logger

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config defines the config for the request-logging middleware.
type Config struct {
	Next        func(c *fiber.Ctx) bool
	SkipBody    func(c *fiber.Ctx) bool
	SkipResBody func(c *fiber.Ctx) bool
	GetResBody  func(c *fiber.Ctx) []byte
	SkipURIs    []string
	Logger      *zap.Logger
	Fields      []string
	FieldsFunc  func(c *fiber.Ctx) []zap.Field
	Messages    []string
	Levels      []zapcore.Level
}

var defaultLogger, _ = zap.NewProduction()

var ConfigDefault = Config{
	Logger:   defaultLogger,
	Fields:   []string{"latency", "status", "method", "url"},
	Messages: []string{"Server error", "Client error", "Success"},
	Levels:   []zapcore.Level{zapcore.ErrorLevel, zapcore.WarnLevel, zapcore.InfoLevel},
}

func configDefault(config ...Config) Config {
	if len(config) < 1 {
		return ConfigDefault
	}
	cfg := config[0]
	if cfg.Logger == nil {
		cfg.Logger = ConfigDefault.Logger
	}
	if cfg.Fields == nil {
		cfg.Fields = ConfigDefault.Fields
	}
	if cfg.Messages == nil {
		cfg.Messages = ConfigDefault.Messages
	}
	if cfg.Levels == nil {
		cfg.Levels = ConfigDefault.Levels
	}
	return cfg
}
