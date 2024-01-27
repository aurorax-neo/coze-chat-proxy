package logger

import (
	"coze-chat-proxy/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// go get -u go.uber.org/zap

// Logger 全局日志对象
var Logger *zap.Logger

func init() {
	// 配置日志输出到文件
	zapConfig := zap.NewProductionConfig()
	zapConfig.OutputPaths = []string{"log.log", "stdout"} // 将日志输出到文件 和 标准输出
	zapConfig.Encoding = "console"                        // 设置日志格 json console
	var LevelErr error
	zapConfig.Level, LevelErr = zap.ParseAtomicLevel(config.CONFIG.LogLevel) // 设置日志级别
	if LevelErr != nil {
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	zapConfig.EncoderConfig = zapcore.EncoderConfig{ // 创建Encoder配置
		MessageKey:   "message",
		LevelKey:     "level",
		EncodeLevel:  zapcore.LowercaseLevelEncoder,
		TimeKey:      "time",
		EncodeTime:   zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		CallerKey:    "caller",
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
	//zapConfig.Sampling = nil

	// 创建Logger对象
	var buildErr error
	Logger, buildErr = zapConfig.Build()
	if buildErr != nil {
		panic("Failed to initialize logger: " + LevelErr.Error())
	}
	// 在应用程序退出时调用以确保所有日志消息都被写入文件
	defer func(Logger *zap.Logger) {
		_ = Logger.Sync()
	}(Logger)
}
