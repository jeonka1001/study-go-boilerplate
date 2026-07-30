package logger

import (
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	PID        string = "pid"
	IP         string = "ip"
	IPs        string = "ips"
	URL        string = "url"
	Latency    string = "latency"
	Status     string = "status"
	Body       string = "body"
	Method     string = "method"
	RequestID  string = "requestId"
	Error      string = "error"
	Protocol   string = "protocol"
	ReqHeaders string = "reqHeaders"
)

// NewMiddleware returns a Fiber request-logging middleware backed by zap.
func NewMiddleware(config ...Config) fiber.Handler {
	cfg := configDefault(config...)
	pid := strconv.Itoa(os.Getpid())

	skipURIs := make(map[string]struct{})
	for _, uri := range cfg.SkipURIs {
		skipURIs[uri] = struct{}{}
	}

	return func(c *fiber.Ctx) error {
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}
		if _, ok := skipURIs[c.Path()]; ok {
			return c.Next()
		}

		start := time.Now()
		chainErr := c.Next()
		latency := time.Since(start)

		status := c.Response().StatusCode()
		index := 2
		switch {
		case status >= 500:
			index = 0
		case status >= 400:
			index = 1
		}
		levelIndex := index
		if levelIndex >= len(cfg.Levels) {
			levelIndex = len(cfg.Levels) - 1
		}
		messageIndex := index
		if messageIndex >= len(cfg.Messages) {
			messageIndex = len(cfg.Messages) - 1
		}

		ce := cfg.Logger.Check(cfg.Levels[levelIndex], cfg.Messages[messageIndex])
		if ce == nil {
			return chainErr
		}

		fields := make([]zap.Field, 0, len(cfg.Fields)+1)
		if cfg.FieldsFunc != nil {
			fields = append(fields, cfg.FieldsFunc(c)...)
		}

		for _, field := range cfg.Fields {
			switch field {
			case "pid":
				fields = append(fields, zap.String(PID, pid))
			case "ip":
				fields = append(fields, zap.String(IP, c.IP()))
			case "ips":
				fields = append(fields, zap.String(IPs, c.Get(fiber.HeaderXForwardedFor)))
			case "url":
				fields = append(fields, zap.String(URL, c.OriginalURL()))
			case "latency":
				fields = append(fields, zap.String(Latency, latency.String()))
			case "status":
				fields = append(fields, zap.Int(Status, status))
			case "body":
				if cfg.SkipBody == nil || !cfg.SkipBody(c) {
					fields = append(fields, zap.ByteString(Body, c.Body()))
				}
			case "method":
				fields = append(fields, zap.String(Method, c.Method()))
			case "requestId":
				fields = append(fields, zap.String(RequestID, c.GetRespHeader(fiber.HeaderXRequestID)))
			case "protocol":
				fields = append(fields, zap.String(Protocol, c.Protocol()))
			case "reqHeaders":
				fields = append(fields, zap.String(ReqHeaders, c.Request().Header.String()))
			case "error":
				if chainErr != nil {
					fields = append(fields, zap.String(Error, chainErr.Error()))
				}
			}
		}

		ce.Write(fields...)
		return chainErr
	}
}
