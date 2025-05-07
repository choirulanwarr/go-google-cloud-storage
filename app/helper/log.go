package helper

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var Logger = NewLoggerLog()
var LogFile *os.File

func NewLoggerLog() *logrus.Logger {
	config := viper.New()

	config.AutomaticEnv()
	config.AddConfigPath("./")
	config.SetConfigName(".env")
	config.SetConfigType("env")

	if err := config.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error load env file: %w", err))
	}

	var writers []io.Writer

	logToStdout := config.GetBool("LOGGER_STDOUT")
	logFileLocation := config.GetString("LOGGER_FILE_LOCATION")
	logLevelStr := config.GetString("LOGGER_LEVEL")

	loggerConf := logrus.New()

	loggerConf.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: "timestamp",
			logrus.FieldKeyMsg:  "message",
		},
	})

	if logFileLocation != "" {
		logLocation := "logs/" + logFileLocation

		logDir := filepath.Dir(logLocation)

		err := os.MkdirAll(logDir, 0755)
		if err != nil {
			loggerConf.Fatalf("Failed to create log directory: %v", err)
		}

		LogFile, err = os.OpenFile(logLocation, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			loggerConf.Fatalf("Failed to open log file: %v", err)
		}

		writers = append(writers, LogFile)
	}

	if logToStdout {
		writers = append(writers, os.Stdout)
	}

	if len(writers) > 0 {
		loggerConf.SetOutput(io.MultiWriter(writers...))
	}

	logLevel, err := logrus.ParseLevel(logLevelStr)
	if err != nil {
		logLevel = logrus.InfoLevel // Default to INFO if parsing fails
	}
	loggerConf.SetLevel(logLevel)

	return loggerConf
}

func CloseLoggerFile() {
	err := LogFile.Close()
	if err != nil {
		Logger.Fatalf("Unable close log file error : %v", err)
	}
}

func LogDebug(apiCallID string, message any) {
	Logger.WithFields(logrus.Fields{"api_call_id": apiCallID}).Debug(message)
}

func LogInfo(apiCallID string, message any) {
	Logger.WithFields(logrus.Fields{"api_call_id": apiCallID}).Info(message)
}

func LogWarning(apiCallID string, message any) {
	Logger.WithFields(logrus.Fields{"api_call_id": apiCallID}).Warn(message)
}

func LogError(apiCallID string, message any) {
	Logger.WithFields(logrus.Fields{"api_call_id": apiCallID}).Error(message)
}
