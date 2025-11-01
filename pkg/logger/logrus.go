package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var logRus *logrus.Logger

func InitLogrus(env string) {
	logRus = logrus.New()

	logRus.SetOutput(os.Stdout)

	if env == "production" {
		logRus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
		})
	} else {
		logRus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     true,
		})
	}

	if env == "development" {
		logRus.SetLevel(logrus.DebugLevel)
	} else {
		logRus.SetLevel(logrus.InfoLevel)
	}

}

func WithFields(fields map[string]interface{}) *logrus.Entry {
	return logRus.WithFields(fields)
}

func InfoLogrus(msg string, keysAndValues ...interface{}) {
	if len(keysAndValues) == 0 {
		logRus.Info(msg)
		return
	}

	logRus.WithFields(toFields(keysAndValues...)).Info(msg)
}

// Debug logs debug level message
func DebugLogrus(msg string, keysAndValues ...interface{}) {
	if len(keysAndValues) == 0 {
		logRus.Debug(msg)
		return
	}
	logRus.WithFields(toFields(keysAndValues...)).Debug(msg)
}

// Warn logs warning level message
func WarnLogrus(msg string, keysAndValues ...interface{}) {
	if len(keysAndValues) == 0 {
		logRus.Warn(msg)
		return
	}
	logRus.WithFields(toFields(keysAndValues...)).Warn(msg)
}

// Error logs error level message
func ErrorLogrus(msg string, keysAndValues ...interface{}) {
	if len(keysAndValues) == 0 {
		logRus.Error(msg)
		return
	}
	logRus.WithFields(toFields(keysAndValues...)).Error(msg)
}

// Fatal logs fatal level message and exits
func FatalLogrus(msg string, keysAndValues ...interface{}) {
	if len(keysAndValues) == 0 {
		logRus.Fatal(msg)
		return
	}
	logRus.WithFields(toFields(keysAndValues...)).Fatal(msg)
}

// Panic logs panic level message and panics
func Panic(msg string, keysAndValues ...interface{}) {
	if len(keysAndValues) == 0 {
		logRus.Panic(msg)
		return
	}
	logRus.WithFields(toFields(keysAndValues...)).Panic(msg)
}

// toFields converts key-value pairs to logrus.Fields
func toFields(keysAndValues ...interface{}) logrus.Fields {
	fields := make(logrus.Fields)
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 < len(keysAndValues) {
			key := keysAndValues[i].(string)
			fields[key] = keysAndValues[i+1]
		}
	}
	return fields
}

func GetLogger() *logrus.Logger {
	return logRus
}
