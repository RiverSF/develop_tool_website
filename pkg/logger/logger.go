package logger

import (
	"log"
	"os"

	"develop_tools/pkg/path"
)

var (
	logsDir = "logs"

	loggerInfo  *log.Logger
	loggerError *log.Logger
)

func Init() error {
	logDir := path.Join(logsDir)
	if err := os.MkdirAll(logDir, 0750); err != nil {
		return err
	}

	var err error
	if loggerInfo, err = newLogger("info"); err != nil {
		return err
	}
	if loggerError, err = newLogger("error"); err != nil {
		return err
	}
	return nil
}

func newLogger(fileName string) (*log.Logger, error) {
	filePath := path.Join(logsDir, fileName+".log")
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0640)
	if err != nil {
		return nil, err
	}

	logger := log.New(file, "", log.Lshortfile|log.Ldate|log.Ltime)

	return logger, nil
}

func Info(f string, v ...interface{}) {
	//标准输出
	log.Printf(f, v...)

	//日志记录
	loggerInfo.Printf(f, v...)
}

func Error(f string, v ...interface{}) {
	//标准输出
	log.Printf(f, v...)

	//日志记录
	loggerError.Printf(f, v...)
}
