package zaputil

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func New(cfg zap.Config) (log *zap.Logger, err error) {
	if len(cfg.EncoderConfig.MessageKey) == 0 || len(cfg.EncoderConfig.LevelKey) == 0 { // 没有正确配置
		return NewDefault(true), nil
	}
	cfg.EncoderConfig.EncodeTime = timeEncoder // 定制时间格式
	if cfg.EncoderConfig.EncodeCaller == nil { // 修正
		cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	}
	if cfg.EncoderConfig.EncodeDuration == nil {
		cfg.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder // 修正
	}
	log, err = cfg.Build()
	return
}

func NewDefault(debugMode bool) *zap.Logger {
	var cfg zap.Config
	if debugMode {
		cfg = zap.NewDevelopmentConfig()
		cfg.OutputPaths = append(cfg.OutputPaths, "stdout.log")
		cfg.ErrorOutputPaths = append(cfg.ErrorOutputPaths, "stderr.log")
	} else {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.EncodeCaller = nil
		cfg.OutputPaths = []string{"stdout"}
	}
	cfg.Encoding = "console"
	cfg.EncoderConfig.EncodeTime = timeEncoder
	log, _ := cfg.Build()
	return log
}

func NewDefaultWithFile(debugMode bool, logFile string) *zap.Logger {
	var cfg zap.Config
	if debugMode {
		cfg = zap.NewDevelopmentConfig()
		cfg.OutputPaths = append(cfg.OutputPaths, logFile)
		cfg.ErrorOutputPaths = append(cfg.ErrorOutputPaths, logFile)
	} else {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.EncodeCaller = nil
	}
	cfg.Encoding = "console"
	cfg.EncoderConfig.EncodeTime = timeEncoder
	log, _ := cfg.Build()
	return log
}
