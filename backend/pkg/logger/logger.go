package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap.Logger for structured logging.
// All logs are output in JSON format for easy parsing by log aggregation systems (Loki, ELK, etc.)
type Logger struct {
	zap *zap.Logger
}

// Config holds logger configuration.
type Config struct {
	Level  string // debug, info, warn, error (default: info)
	Format string // json, text (default: json)
	Output string // stdout, stderr, or file path (default: stdout)
}

// New creates a new structured logger with JSON output (suitable for Grafana/Loki).
func New(cfg Config) (*Logger, error) {
	// Parse log level
	level := parseLevel(cfg.Level)

	// Configure encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.MessageKey = "message"
	encoderConfig.LevelKey = "level"
	encoderConfig.CallerKey = "caller"

	// Choose encoder based on format
	var encoder zapcore.Encoder
	if cfg.Format == "text" {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// Choose output
	var writeSyncer zapcore.WriteSyncer
	switch cfg.Output {
	case "stderr":
		writeSyncer = zapcore.AddSync(os.Stderr)
	case "stdout", "":
		writeSyncer = zapcore.AddSync(os.Stdout)
	default:
		// File output
		file, err := os.OpenFile(cfg.Output, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		writeSyncer = zapcore.AddSync(file)
	}

	core := zapcore.NewCore(encoder, writeSyncer, level)
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &Logger{zap: zapLogger}, nil
}

// NewDefault creates a new logger with default JSON output configuration.
func NewDefault() *Logger {
	cfg := Config{
		Level:  "info",
		Format: "json",
		Output: "stdout",
	}
	logger, _ := New(cfg)
	return logger
}

// parseLevel converts string level to zapcore.Level.
func parseLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info", "":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// With creates a child logger with additional fields.
func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{zap: l.zap.With(fields...)}
}

// WithRequestID adds request ID to logger context.
func (l *Logger) WithRequestID(requestID string) *Logger {
	return l.With(zap.String("request_id", requestID))
}

// Info logs an info message with optional fields.
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.zap.Info(msg, fields...)
}

// Error logs an error message with optional fields.
func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.zap.Error(msg, fields...)
}

// Debug logs a debug message with optional fields.
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.zap.Debug(msg, fields...)
}

// Warn logs a warning message with optional fields.
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.zap.Warn(msg, fields...)
}

// Fatal logs a fatal message and exits.
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.zap.Fatal(msg, fields...)
}

// Sync flushes any buffered log entries.
func (l *Logger) Sync() error {
	return l.zap.Sync()
}

// Zap returns the underlying zap logger for advanced usage.
func (l *Logger) Zap() *zap.Logger {
	return l.zap
}

// Helper functions for common field types
func String(key, val string) zap.Field {
	return zap.String(key, val)
}

func Int(key string, val int) zap.Field {
	return zap.Int(key, val)
}

func Int64(key string, val int64) zap.Field {
	return zap.Int64(key, val)
}

func Float64(key string, val float64) zap.Field {
	return zap.Float64(key, val)
}

func Bool(key string, val bool) zap.Field {
	return zap.Bool(key, val)
}

func Error(err error) zap.Field {
	return zap.Error(err)
}

func Duration(key string, d time.Duration) zap.Field {
	return zap.Duration(key, d)
}
