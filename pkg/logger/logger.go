package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"layout-for-go/pkg/env"
	"os"
	"sync"
)

// option 日志配置结构体
type option struct {
	level  zapcore.Level
	fields map[string]string
}

// Option 日志配置函数变量
type Option func(*option)

// WithLevel 配置日志等级
func WithLevel(level zapcore.Level) Option {
	return func(o *option) {
		o.level = level
	}
}

// WithField 配置额外固定字段到日志中
func WithField(key, value string) Option {
	return func(o *option) {
		if key != "" && value != "" {
			o.fields[key] = value
		}
	}
}

// 单例获取 logger
var (
	once     sync.Once
	instance *zap.Logger
)

// Get 获取 logger 对象
func Get(opts ...Option) zap.Logger {
	once.Do(func() {
		// 初始化默认配置
		opt := &option{
			level:  zapcore.InfoLevel,
			fields: make(map[string]string),
		}

		// 应用传入配置到默认 opt 中
		for _, apply := range opts {
			apply(opt)
		}

		cfg := zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}

		// dev 环境使用 console 输出，其他环境选择 json 输出
		var encoder zapcore.Encoder
		if env.Current() == env.Dev {
			encoder = zapcore.NewConsoleEncoder(cfg)
		} else {
			encoder = zapcore.NewJSONEncoder(cfg)
		}

		core := zapcore.NewCore(
			encoder,
			zapcore.Lock(os.Stdout),
			opt.level,
		)
		instance = zap.New(
			core,
			zap.AddCaller(),
			zap.AddStacktrace(zap.ErrorLevel),
		)

		// 添加额外字段到 logger 中
		fields := make([]zap.Field, 0, len(opt.fields))
		for key, value := range opt.fields {
			if key == "" || value == "" {
				continue
			}
			fields = append(fields, zap.String(key, value))
		}
		if len(fields) > 0 {
			instance = instance.WithOptions(zap.Fields(fields...))
		}
	})

	return *instance
}
