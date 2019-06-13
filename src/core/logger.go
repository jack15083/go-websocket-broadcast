package core

import (
	"os"
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
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

	maxAge := Config.Logger.MaxAge * time.Hour
	rotationTime := Config.Logger.RotationTime * time.Hour

	l.ConfigLocalFilesystemLogger(Config.Logger.LogPath, Config.AppName, maxAge, rotationTime)
}

// config logrus log to local filesystem, with file rotation
func (l *Logger) ConfigLocalFilesystemLogger(logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration) {
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
	}, &log.TextFormatter{})

	log.AddHook(lfHook)
}
