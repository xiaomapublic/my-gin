//日志服务类，使用zap库
package log

import (
	"my-gin/libraries/config"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLog(logName string) *zap.SugaredLogger {
	//config.DefaultConfigInit()
	logConfigs := config.UnmarshalConfig.Log

	hook := lumberjack.Logger{

		Filename:   logConfigs.Path + time.Now().Format("2006-01-02") + "/" + logName + ".log", // 日志文件路径
		MaxSize:    logConfigs.Max_size,                                                        // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: logConfigs.Max_backups,                                                     // 日志文件最多保存多少个备份
		MaxAge:     logConfigs.Max_age,                                                         // 文件最多保存多少天
		Compress:   logConfigs.Compress,                                                        // 是否压缩
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "log",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.InfoLevel)

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),            // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(&hook)), // 打印到文件
		atomicLevel, // 日志级别
	)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 构造日志
	var sugar *zap.SugaredLogger

	sugar = zap.New(core, caller, development).Sugar()

	defer sugar.Sync()

	return sugar

}
