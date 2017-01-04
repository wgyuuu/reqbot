package dlog

import (
	"log/syslog"
	"os"
	"time"

	"fmt"

	"github.com/Sirupsen/logrus"
)

type Entry struct {
	*logrus.Entry
}

func NewEntry(fields map[string]interface{}) *Entry {
	return &Entry{log.WithFields(fields)}
}

func (e *Entry) Info2(args ...interface{}) {
	e.Info(formatKV(args...))
}

func (e *Entry) Debug2(args ...interface{}) {
	e.Debug(formatKV(args...))
}

var log = logrus.New()

func Init(isDebug bool) {
	log.Out = os.Stderr
	log.Formatter = &logrus.TextFormatter{
		DisableColors:   true,
		TimestampFormat: time.RFC1123,
	}
	if isDebug {
		log.Level = logrus.DebugLevel
	}

	hook, err := NewSyslogHook("udp", "127.0.0.1:8000", syslog.LOG_DEBUG, "")
	if err != nil {
		logrus.Fatal(err)
		return
	}
	log.Hooks.Add(hook)
}

func Info(args ...interface{}) {
	log.Info(formatKV(args...))
}

func Error(args ...interface{}) {
	log.Error(formatKV(args...))
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func formatKV(kv ...interface{}) string {
	if len(kv) == 0 {
		return ""
	}
	if len(kv)%2 != 0 {
		kv = append(kv, "unknow")
	}

	format := ""
	for i := 0; i < len(kv); i += 2 {
		format += "||%v=%+v"
	}
	return fmt.Sprintf(format[2:], kv...)
}
