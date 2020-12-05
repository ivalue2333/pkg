package logx

import "github.com/sirupsen/logrus"

func newJSONFormatter() logrus.Formatter {
	formatter := new(JSONFormatter)
	formatter.TimestampFormat = "2006-01-02T15:04:05.000Z07:00"
	return formatter
}

func newTextFormatter() logrus.Formatter {
	formatter := new(logrus.TextFormatter)
	formatter.TimestampFormat = "2006-01-02T15:04:05.000Z07:00"
	return formatter
}
