package util

import (
	"fmt"
	"log"
	"os"
)

var Log *logger

type logger struct {
	info  *log.Logger
	warn  *log.Logger
	err   *log.Logger
	fatal *log.Logger
}

const (
	infoPrefix  = "[INFO]: "
	warnPrefix  = "[WARN]: "
	errorPrefix = "[ERROR]: "
)

func InitLog() {
	Log = &logger{
		info:  log.New(os.Stdout, infoPrefix, log.Ldate|log.Ltime),
		warn:  log.New(os.Stdout, warnPrefix, log.Ldate|log.Ltime|log.Llongfile),
		err:   log.New(os.Stdout, errorPrefix, log.Ldate|log.Ltime|log.Llongfile),
		fatal: log.New(os.Stdout, errorPrefix, log.Ldate|log.Ltime|log.Llongfile),
	}

}

func (l *logger) Info(msg string, v ...interface{}) {
	msg = msg + "\n"
	_ = l.info.Output(2, fmt.Sprintf(msg, v...))
}

func (l *logger) Warn(msg string, v ...interface{}) {
	msg = msg + "\n"
	_ = l.warn.Output(2, fmt.Sprintf(msg, v...))
}

func (l *logger) Error(msg string, v ...interface{}) {
	msg = msg + "\n"
	_ = l.err.Output(2, fmt.Sprintf(msg, v...))
}

func (l *logger) Fatal(msg string, v ...interface{}) {
	msg = msg + "\n"
	_ = l.err.Output(2, fmt.Sprintf(msg, v...))
	os.Exit(1)
}
