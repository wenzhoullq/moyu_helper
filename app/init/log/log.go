package log

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"time"
	"weixin_LLM/init/config"
)

var Logger *logrus.Logger

func InitLog() error {

	logFilePath := config.Config.LogFilePath
	logFileName := config.Config.LogFileName

	// 日志文件
	fileName := path.Join(logFilePath, logFileName)
	//检查日志是否存在
	f := &os.File{}
	_, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(filepath.Dir(fileName), os.ModePerm); err != nil {
				return err
			}
			f, err = os.Create(fileName)
			if err != nil {
				return err
			}
			f.Close()
		}
	}
	// 写入文件
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}

	// 实例化
	Logger = logrus.New()

	// 设置输出
	Logger.Out = src

	// 设置日志级别
	Logger.SetLevel(logrus.DebugLevel)

	// 设置 rotatelogs
	logWriter, err := rotatelogs.New(
		// 分割后的文件名称
		fileName+".%Y%m%d.log",

		// 设置最大保存时间(7天)
		rotatelogs.WithMaxAge(7*24*time.Hour),

		// 设置日志切割时间间隔(1天)
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		return err
	}
	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}

	lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 新增 Hook
	Logger.AddHook(lfHook)
	return nil
}
