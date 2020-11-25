package logx

import "github.com/sirupsen/logrus"

type Entry = logrus.Entry
type Ext1FieldLogger = logrus.Ext1FieldLogger
type FieldLogger = logrus.FieldLogger
type Fields = logrus.Fields
type Formatter = logrus.Formatter
type Hook = logrus.Hook
type Level = logrus.Level
type LevelHooks = logrus.LevelHooks
type MutexWrap = logrus.MutexWrap

const PanicLevel = logrus.PanicLevel
const FatalLevel = logrus.FatalLevel
const ErrorLevel = logrus.ErrorLevel
const WarnLevel = logrus.WarnLevel
const InfoLevel = logrus.InfoLevel
const DebugLevel = logrus.DebugLevel
const TraceLevel = logrus.TraceLevel


var (
	AllLevels = logrus.AllLevels
)