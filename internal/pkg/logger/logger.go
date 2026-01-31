package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger 全局日志实例
var log *zap.Logger

// Config 日志配置
type Config struct {
	Level  string `mapstructure:"level"`  // debug, info, warn, error
	Format string `mapstructure:"format"` // json, console
	Output string `mapstructure:"output"` // stdout, file, both
	Path   string `mapstructure:"path"`   // 日志文件路径
}

// parseLevel 解析日志级别配置
func parseLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// buildEncoderConfig 构建编码器配置
func buildEncoderConfig(format string) zapcore.EncoderConfig {
	common := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	if format == "console" {
		common.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		common.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	return common
}

// buildFileCore 构建文件输出 Core
func buildFileCore(encoderConfig zapcore.EncoderConfig, path string, level zapcore.Level) ([]zapcore.Core, error) {
	if path == "" {
		path = "logs"
	}
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	var cores []zapcore.Core
	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)

	// Info 日志
	infoPath := filepath.Join(path, "info.log")
	infoFile, err := os.OpenFile(infoPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open info log file: %w", err)
	}
	defer infoFile.Close()
	cores = append(cores, zapcore.NewCore(jsonEncoder, zapcore.AddSync(infoFile), level))

	// Error 日志
	errorPath := filepath.Join(path, "error.log")
	errorFile, err := os.OpenFile(errorPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open error log file: %w", err)
	}
	defer errorFile.Close()
	cores = append(cores, zapcore.NewCore(jsonEncoder, zapcore.AddSync(errorFile), zapcore.ErrorLevel))

	return cores, nil
}

// Init 初始化日志模块
func Init(cfg *Config) error {
	level := parseLevel(cfg.Level)
	encoderConfig := buildEncoderConfig(cfg.Format)

	var cores []zapcore.Core

	// 控制台输出
	if cfg.Output == "stdout" || cfg.Output == "both" {
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level))
	}

	// 文件输出
	if cfg.Output == "file" || cfg.Output == "both" {
		fileCores, err := buildFileCore(encoderConfig, cfg.Path, level)
		if err != nil {
			return err
		}
		cores = append(cores, fileCores...)
	}

	log = zap.New(zapcore.NewTee(cores...),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	return nil
}

// Get 获取日志实例
func Get() *zap.Logger {
	if log == nil {
		// 默认配置
		_ = Init(&Config{
			Level:  "info",
			Format: "json",
			Output: "stdout",
		})
	}
	return log
}

// Sync 刷新日志
func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}

// Debug 调试日志
// TODO: 确认是否需要此函数
func Debug(msg string, fields ...zap.Field) {
	Get().Debug(msg, fields...)
}

// Info 信息日志
func Info(msg string, fields ...zap.Field) {
	Get().Info(msg, fields...)
}

// Warn 警告日志
// TODO: 确认是否需要此函数
func Warn(msg string, fields ...zap.Field) {
	Get().Warn(msg, fields...)
}

// Error 错误日志
func Error(msg string, fields ...zap.Field) {
	Get().Error(msg, fields...)
}

// Fatal 致命日志
func Fatal(msg string, fields ...zap.Field) {
	Get().Fatal(msg, fields...)
}

// ===== Zap Field 辅助函数 =====

// String 创建 string 类型字段
func String(key string, val string) zap.Field {
	return zap.String(key, val)
}

// Int 创建 int 类型字段
func Int(key string, val int) zap.Field {
	return zap.Int(key, val)
}

// Any 创建任意类型字段
func Any(key string, val interface{}) zap.Field {
	return zap.Any(key, val)
}

// Err 创建 error 类型字段
func Err(err error) zap.Field {
	return zap.Error(err)
}

// Duration 创建 duration 类型字段
func Duration(key string, val interface{}) zap.Field {
	return zap.Duration(key, val.(time.Duration))
}
