package log

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

// error logger
var myLogger *zap.SugaredLogger

var levelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func getLoggerLevel(lvl string) zapcore.Level {
	if level, ok := levelMap[lvl]; ok {
		return level
	}
	return zapcore.InfoLevel
}

func NewLog(level string, isFile bool) *zap.SugaredLogger {

	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(getLoggerLevel(level))
	var core zapcore.Core
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,                       //控制台彩色日志输出
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"), //时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,                         // 时间精度？
		EncodeCaller:   zapcore.ShortCallerEncoder,                             // 短路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	if isFile {
		softDir := "."
		if dir, err := os.Getwd(); err != nil {
			softDir = dir
		}
		fileName := fmt.Sprintf("%s/logs/%s.log", softDir, time.Now().Format("20060102"))
		fmt.Println("当前日志文件：", fileName)
		hook := lumberjack.Logger{
			Filename:   fileName, // 日志文件路径
			MaxSize:    128,      // 每个日志文件保存的最大尺寸 单位：M
			MaxBackups: 30,       // 日志文件最多保存多少个备份
			MaxAge:     7,        // 文件最多保存多少天
			Compress:   true,     // 是否压缩
		}
		encoderConfig.EncodeLevel=zapcore.CapitalLevelEncoder
		core = zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),                                        // 日志格式
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
			atomicLevel,                                                                     // 日志级别
		)

	} else {
		core = zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),                // 日志格式
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), // 打印到控制台和文件
			atomicLevel,                                             // 日志级别
		)
	}
	//日志级别=debug时，
	if level == "debug" {
		caller := zap.AddCaller()           //开启开发模式，堆栈跟踪
		development := zap.AddCallerSkip(0) //开启文件及行号
		logger := zap.New(core, caller, development)
		myLogger = logger.Sugar()
	} else {
		logger := zap.New(core)
		myLogger = logger.Sugar()
	}
	return myLogger
}


//兼容 log.Println [INFO]级别
func Println(args ...interface{}) {
	myLogger.Info(args...)
}

//兼容 log.Printf [INFO]级别
func Printf(template string, args ...interface{}) {
	myLogger.Infof(template, args...)
}

func Debug(args ...interface{}) {
	myLogger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	myLogger.Debugf(template, args...)
}

func Info(args ...interface{}) {
	myLogger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	myLogger.Infof(template, args...)
}

func Warn(args ...interface{}) {
	myLogger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	myLogger.Warnf(template, args...)
}

func Error(args ...interface{}) {
	myLogger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	myLogger.Errorf(template, args...)
}

func DPanic(args ...interface{}) {
	myLogger.DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	myLogger.DPanicf(template, args...)
}

func Panic(args ...interface{}) {
	myLogger.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	myLogger.Panicf(template, args...)
}

func Fatal(args ...interface{}) {
	myLogger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	myLogger.Fatalf(template, args...)
}
