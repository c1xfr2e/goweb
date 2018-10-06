package logging

import (
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zaoshu/hardcore/logging/fluentd"
)

// Config log config
type Config struct {
	Level     string          `json:"level" yaml:"level"`         // `debug`, `info`, `warn`, `error`
	Formatter string          `json:"formatter" yaml:"formatter"` // ["json", "text"], default is json
	Fluent    *fluentd.Config `json:"fluent" yaml:"fluent"`       // 设置后会同时将日志输出到 stdout 和 fluent
}

// InitFromConfig init log from config
func InitFromConfig(c Config) error {
	logrus.SetFormatter(getFormatter(c.Formatter))

	var l logrus.Level
	if len(c.Level) > 0 {
		var err error
		l, err = logrus.ParseLevel(c.Level)
		if err != nil {
			return err
		}
	} else {
		l = logrus.InfoLevel
	}
	logrus.SetLevel(l)
	logrus.Infof("log level %s", l)

	if c.Fluent == nil {
		return nil
	}

	tags := strings.Split(c.Fluent.Tag, ".")
	if len(tags) != 2 {
		return fmt.Errorf("invalid fluent tag %s, tag pattern is xxx.xxx", c.Fluent.Tag)
	}
	switch tags[0] {
	case "dev":
	case "debug":
	case "release":
	default:
		return fmt.Errorf("invalid fluent tag %s, first part of tag must be one of dev, debug, release", c.Fluent.Tag)
	}

	hook, err := fluentd.NewFluentHook(*c.Fluent)
	if err != nil {
		return err
	}
	logrus.AddHook(hook)
	return nil
}

func getFormatter(formatter string) logrus.Formatter {
	switch formatter {
	case "":
		fallthrough
	case "json":
		return &logrus.JSONFormatter{TimestampFormat: time.RFC3339Nano}
	case "text":
		return &TextFormatter{FullTimestamp: true}
	}
	panic(fmt.Errorf("invalid log formatter %s", formatter))
}

// NewLogger new logger
//
// fields: set logger fields, can be nil
func NewLogger(fields map[string]interface{}) *logrus.Entry {
	if len(fields) > 0 {
		return logrus.WithFields(fields)
	}
	return logrus.NewEntry(logrus.StandardLogger())
}

// NewLoggerWithRequestID new logger and add field requestID
func NewLoggerWithRequestID(requestID string) *logrus.Entry {
	if len(requestID) > 0 {
		return logrus.WithField("requestID", requestID)
	}
	return logrus.NewEntry(logrus.StandardLogger())
}
