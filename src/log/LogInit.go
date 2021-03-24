package log

import (
	"SHUCTES/src/config"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
)


var Logger = logrus.New()

func InitLog() {
	logPath := path.Join(config.Conf.LogCon.FilePath, config.Conf.LogCon.FileName)

	//0666 -rw-rw-rw- 创建文件所有者，用户组和其他人对该文件有读写权限
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	//同时写文件和console
	writers := []io.Writer{
		file,
		os.Stdout}
	fileAndStdoutWriter := io.MultiWriter(writers...)

	if err != nil {
		Logger.Error("failed to log to file.")
	} else {
		Logger.SetOutput(fileAndStdoutWriter)
	}

	//设置日期格式
	Logger.SetFormatter(&logrus.TextFormatter{
	TimestampFormat:
		"2006-01-02 15:04:05",
	})

	//设置最低loglevel
	Logger.SetLevel(logrus.DebugLevel)

	Logger.Info("LogConf init successfully")
	//Panic, Fatal, Error, Warn, Info, Debug, Trace
}
