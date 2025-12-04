package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	infoLogger  = log.New(os.Stdout, "", 0)
	errorLogger = log.New(os.Stderr, "", 0)
)

func ts() string {
	return time.Now().Format("02:01:2006 15:04:05")
}

func Infof(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	infoLogger.Printf("[%s] [INFO] %s", ts(), msg)
}

func Errorf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	errorLogger.Printf("[%s] [ERROR] %s", ts(), msg)
}
