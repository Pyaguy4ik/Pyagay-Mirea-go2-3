package logger

import (
    "os"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

// Config конфигурация логгера
type Config struct {
    ServiceName string
    Environment string // development/production
    LogLevel    string
}

// New создает новый экземпляр логгера
func New(cfg Config) (*zap.Logger, error) {
    // Определяем уровень логирования
    var level zapcore.Level
    switch cfg.LogLevel {
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

    // Настройка энкодера
    encoderConfig := zapcore.EncoderConfig{
        TimeKey:        "ts",
        LevelKey:       "level",
        NameKey:        "logger",
        CallerKey:      "caller",
        MessageKey:     "msg",
        StacktraceKey:  "stacktrace",
        LineEnding:     zapcore.DefaultLineEnding,
        EncodeLevel:    zapcore.LowercaseLevelEncoder,
        EncodeTime:     zapcore.ISO8601TimeEncoder,
        EncodeDuration: zapcore.MillisDurationEncoder,
        EncodeCaller:   zapcore.ShortCallerEncoder,
    }

    // Выбираем формат: JSON для production, консоль для development
    var encoder zapcore.Encoder
    if cfg.Environment == "production" {
        encoder = zapcore.NewJSONEncoder(encoderConfig)
    } else {
        encoder = zapcore.NewConsoleEncoder(encoderConfig)
    }

    // Создаем core
    core := zapcore.NewCore(
        encoder,
        zapcore.AddSync(os.Stdout),
        level,
    )

    // Добавляем служебные поля
    logger := zap.New(core).With(
        zap.String("service", cfg.ServiceName),
    )

    return logger, nil
}

// WithRequestID добавляет request_id в логгер
func WithRequestID(logger *zap.Logger, requestID string) *zap.Logger {
    return logger.With(zap.String("request_id", requestID))
}

// WithComponent добавляет компонент в логгер
func WithComponent(logger *zap.Logger, component string) *zap.Logger {
    return logger.With(zap.String("component", component))
}
