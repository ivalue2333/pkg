package logx

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
)

func ParseLevel(level string) (Level, error) {
	return logrus.ParseLevel(level)
}

func HandleFileOutput(l Logger, fileName string) error {
	writer, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return errors.Wrapf(err, "failed to open file(%s)", fileName)
	}
	l.SetOutput(writer) // 设置正常日志
	return nil
}
