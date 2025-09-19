package logger

import (
	"fmt"
	"os"
	"time"
	_ "unsafe"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type myLog struct {
	ZapLog   *zap.Logger
	logLevel zapcore.Level
}

func init() {
	InitLog(zapcore.DebugLevel)
}

var Log *myLog

func (log *myLog) Info(v ...interface{}) {
	if log.logLevel > zap.InfoLevel {
		return
	}
	message := fmt.Sprint(v...)
	fields := []zap.Field{}
	if ce := log.ZapLog.Check(zap.InfoLevel, message); ce != nil {
		ce.Write(fields...)
	}
}

func (log *myLog) Infof(format string, v ...interface{}) {
	if log.logLevel > zap.InfoLevel {
		return
	}
	message := fmt.Sprintf(format, v...)
	fields := []zap.Field{}
	if ce := log.ZapLog.Check(zap.InfoLevel, message); ce != nil {
		ce.Write(fields...)
	}
}

func (log *myLog) Error(v ...interface{}) {
	if log.logLevel > zap.ErrorLevel {
		return
	}
	message := fmt.Sprint(v...)
	fields := []zap.Field{}
	if ce := log.ZapLog.Check(zap.ErrorLevel, message); ce != nil {
		ce.Write(fields...)
	}
}
func (log *myLog) Errorf(format string, v ...interface{}) {
	if log.logLevel > zap.ErrorLevel {
		return
	}
	message := fmt.Sprintf(format, v...)
	fields := []zap.Field{}
	if ce := log.ZapLog.Check(zap.ErrorLevel, message); ce != nil {
		ce.Write(fields...)
	}
}
func (log *myLog) Warn(v ...interface{}) {
	if log.logLevel > zap.WarnLevel {
		return
	}
	message := fmt.Sprint(v...)
	fields := []zap.Field{}
	if ce := log.ZapLog.Check(zap.WarnLevel, message); ce != nil {
		ce.Write(fields...)
	}
}
func (log *myLog) Warnf(format string, v ...interface{}) {
	if log.logLevel > zap.WarnLevel {
		return
	}
	message := fmt.Sprintf(format, v...)
	fields := []zap.Field{}
	if ce := log.ZapLog.Check(zap.WarnLevel, message); ce != nil {
		ce.Write(fields...)
	}
}

func (log *myLog) Debug(v ...interface{}) {
	if log.logLevel > zap.DebugLevel {
		return
	}
	message := fmt.Sprint(v...)
	fields := []zap.Field{}
	if ce := log.ZapLog.Check(zap.DebugLevel, message); ce != nil {
		ce.Write(fields...)
	}
}
func (log *myLog) Debugf(format string, v ...interface{}) {
	if log.logLevel > zap.DebugLevel {
		return
	}
	message := fmt.Sprintf(format, v...)
	fields := []zap.Field{}
	if ce := log.ZapLog.Check(zap.DebugLevel, message); ce != nil {
		ce.Write(fields...)
	}
}
func (log *myLog) Fatal(v ...interface{}) {
	message := fmt.Sprint(v...)
	fields := []zap.Field{}
	if ce := log.ZapLog.Check(zap.FatalLevel, message); ce != nil {
		ce.Write(fields...)
	}
	os.Exit(1)
}
func (log *myLog) Fatalf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	fields := []zap.Field{}
	if ce := log.ZapLog.Check(zap.FatalLevel, message); ce != nil {
		ce.Write(fields...)
	}
	os.Exit(1)
}
func (log *myLog) Panic(v ...interface{}) {
	message := fmt.Sprint(v...)
	fields := []zap.Field{}
	if ce := log.ZapLog.Check(zap.PanicLevel, message); ce != nil {
		ce.Write(fields...)
	}
	panic(message)
}
func (log *myLog) Panicf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	fields := []zap.Field{}
	if ce := log.ZapLog.Check(zap.PanicLevel, message); ce != nil {
		ce.Write(fields...)
	}
	panic(message)
}

var levelMap = map[string]zapcore.Level{

	"debug": zapcore.DebugLevel,

	"info": zapcore.InfoLevel,

	"warn": zapcore.WarnLevel,

	"error": zapcore.ErrorLevel,

	"dpanic": zapcore.DPanicLevel,

	"panic": zapcore.PanicLevel,

	"fatal": zapcore.FatalLevel,
}

// 生成日志文件
func setLoggerFile(lev_name string, encoder zapcore.Encoder, enableLogFile bool) zapcore.Core {
	//日志级别
	priority := zap.LevelEnablerFunc(func(level zapcore.Level) bool { //info和debug级别,debug级别是最低的
		return levelMap[lev_name] == level
	})

	// filename := fmt.Sprintf(`%s/%s-%s-%s/%s.log`, conf.Config.Log.Path, time.Now().Format("2006"), time.Now().Format("01"), time.Now().Format("02"), lev_name)
	filename := fmt.Sprintf(`%s/%s-%s-%s/%s.log`, "./log", time.Now().Format("2006"), time.Now().Format("01"), time.Now().Format("02"), lev_name)
	//info文件writeSyncer
	infoFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filename, //日志文件存放目录，如果文件夹不存在会自动创建
		MaxSize:    10,       //文件大小限制,单位MB
		MaxBackups: 100,      //最大保留日志文件数量
		MaxAge:     30,       //日志文件保留天数
		Compress:   false,    //是否压缩处理
	})
	if enableLogFile {
		return zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(infoFileWriteSyncer, zapcore.AddSync(os.Stdout)), priority)
	}
	return zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), priority)
}

func InitLog(logLevel zapcore.Level) {
	var coreArr []zapcore.Core // 存储核心组件
	m := &myLog{
		ZapLog:   nil,
		logLevel: logLevel,
	}
	//获取编码器:

	encoderConfig := zap.NewProductionEncoderConfig()                     //NewJSONEncoder()输出json格式，NewConsoleEncoder()输出普通文本格式
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime) //指定时间格式
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder          //按级别显示不同颜色，不需要的话取值zapcore.CapitalLevelEncoder就可以了
	//encoderConfig.EncodeCaller = zapcore.FullCallerEncoder        //显示完整文件路径
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	//encoder := zapcore.NewJSONEncoder(encoderConfig)
	_ = encoder
	// zzk: 日志
	for key := range levelMap {
		fileCore := setLoggerFile(key, encoder, false) // zzk: 配置生成在 后端代码的同级目录。 所以这部分代码应该不用改
		coreArr = append(coreArr, fileCore)
	}
	m.ZapLog = zap.New(zapcore.NewTee(coreArr...), zap.AddCaller(), zap.AddCallerSkip(1)) // zzk: 要分成两个选项来添加是不是显得有些多余？不过源码确实这样的

	Log = m
}
