package core

import (
	"os"
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func (l *Logger) Init() {
	if Config.Logger.Debug {
		log.SetLevel(log.DebugLevel)
	}

	if !Config.Logger.OutFile {
		log.SetOutput(os.Stdout)
		return
	}

	os.MkdirAll(Config.Logger.LogPath, 0755)

	maxAge := Config.Logger.MaxAge * time.Hour
	rotationTime := Config.Logger.RotationTime * time.Hour

	log.AddHook(l.NewRotateHook(Config.Logger.LogPath, Config.AppName, maxAge, rotationTime))
}

// config logrus log to local filesystem, with file rotation
func (l *Logger) NewRotateHook(logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration) *lfshook.LfsHook {
	baseLogPaht := path.Join(logPath, logFileName)
	writer, err := rotatelogs.New(
		baseLogPaht+".%Y%m%d.log",
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)

	if err != nil {
		log.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer, // 为不同级别设置不同的输出目的
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
		log.FatalLevel: writer,
		log.PanicLevel: writer,
	}, &log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05.000"})

	return lfHook
}

func (l *Logger) New(logFileName string) *logrus.Logger {
	var lognew = logrus.New()

	if Config.Logger.Debug {
		lognew.SetLevel(log.DebugLevel)
	}

	maxAge := Config.Logger.MaxAge * time.Hour
	rotationTime := Config.Logger.RotationTime * time.Hour

	lognew.AddHook(l.NewRotateHook(Config.Logger.LogPath, logFileName, maxAge, rotationTime))
	return lognew
}
