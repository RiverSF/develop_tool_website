package logger

import (
	"log"
	"os"
	"sync"

	"develop_tools/pkg/path"
)

var (
	logsDir = "logs"

	loggerInfo  *log.Logger
	loggerError *log.Logger

	infoFile  *os.File
	errorFile *os.File
	closeOnce sync.Once
)

func Init() error {
	logDir := path.Join(logsDir)
	if err := os.MkdirAll(logDir, 0750); err != nil {
		return err
	}

	var err error
	if loggerInfo, infoFile, err = newLogger("info"); err != nil {
		return err
	}
	if loggerError, errorFile, err = newLogger("error"); err != nil {
		_ = infoFile.Close()
		return err
	}
	return nil
}

func Close() error {
	var err error
	closeOnce.Do(func() {
		if infoFile != nil {
			if e := infoFile.Close(); e != nil {
				err = e
			}
		}
		if errorFile != nil {
			if e := errorFile.Close(); e != nil && err == nil {
				err = e
			}
		}
	})
	return err
}

func newLogger(fileName string) (*log.Logger, *os.File, error) {
	filePath := path.Join(logsDir, fileName+".log")
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0640)
	if err != nil {
		return nil, nil, err
	}

	return log.New(file, "", log.Lshortfile|log.Ldate|log.Ltime), file, nil
}

func Info(f string, v ...interface{}) {
	log.Printf(f, v...)
	loggerInfo.Printf(f, v...)
}

func Error(f string, v ...interface{}) {
	log.Printf(f, v...)
	loggerError.Printf(f, v...)
}
