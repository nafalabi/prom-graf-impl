package utils

import (
	"log"
	"os"
)

const (
	DEBUG = iota
	INFO
	WARN
	ERROR
)

type logger struct {
	LogLevel    int
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
}

func (l *logger) Debug(message string) {
	if l.LogLevel <= DEBUG {
		l.debugLogger.Println(message)
	}
}

func (l *logger) Info(message string) {
	if l.LogLevel <= INFO {
		l.infoLogger.Println(message)
	}
}

func (l *logger) Warn(message string) {
	if l.LogLevel <= WARN {
		l.warnLogger.Println(message)
	}
}

func (l *logger) Error(message string) {
	l.errorLogger.Println(message)
}

func NewLogger(logLevel int) logger {
	l := logger{}
	l.debugLogger = log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
	l.infoLogger = log.New(os.Stdout, "[INFO]  ", log.Ldate|log.Ltime|log.Lshortfile)
	l.warnLogger = log.New(os.Stdout, "[WARN]  ", log.Ldate|log.Ltime|log.Lshortfile)
	l.errorLogger = log.New(os.Stdout, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	return l
}
