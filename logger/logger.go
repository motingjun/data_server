package logger

import (
	"time"

	"data_server/hooks"

	rotations "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Log struct {
	Logger *logrus.Logger
}

func getLogLevel() logrus.Level {
	level := viper.GetString("LoggingLevel")
	switch level {
	case "TRACE":
		return logrus.TraceLevel
	case "DEBUG":
		return logrus.DebugLevel
	case "INFO":
		return logrus.InfoLevel
	case "ERROR":
		return logrus.ErrorLevel
	case "FATAL":
		return logrus.FatalLevel
	case "PANIC":
		return logrus.FatalLevel
	default:
		return logrus.InfoLevel
	}
}

func InitLog() *Log {

	logger := Log{
		Logger: logrus.New(),
	}
	// logName := viper.GetString("LogFilePrefix") + "." + strconv.Itoa(os.Getpid()) + ".%Y-%m-%d_%H"
	logName := viper.GetString("LogFilePrefix") + "_%Y-%m-%d_%H" + ".log"
	rotater, err := rotations.New(
		logName,
		rotations.WithLinkName(logName),
		rotations.WithRotationTime(time.Duration(viper.GetInt("LogRotationInterval"))*time.Second),
		rotations.WithMaxAge(time.Duration(viper.GetInt("LogMaxAge"))*time.Second),
	)
	if err != nil {
		panic("Create a rotation of log failed")
	}

	logger.Logger.Hooks.Add(hooks.NewContextHook())
	logger.Logger.Formatter = &logrus.TextFormatter{}
	logger.Logger.SetLevel(getLogLevel())
	logger.Logger.Out = rotater

	return &logger
}
