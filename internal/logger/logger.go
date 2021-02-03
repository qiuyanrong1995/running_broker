package logger

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"strings"
)

func NewLogger(level string, rotation int, size int, age int, path string, name string, json bool, debug bool) (*logrus.Logger, error) {
	logPath := getLogPath(path, name)
	var writers []io.Writer
	rotationWriter := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    size,
		MaxAge:     age,
		MaxBackups: rotation,
	}
	writers = append(writers, rotationWriter)

	logger := logrus.New()
	if json {
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{})
	}

	var logLevel logrus.Level
	var err error

	if debug {
		writers = append(writers, os.Stdout)
		logLevel, err = logrus.ParseLevel("debug")
	} else {
		logLevel, err = logrus.ParseLevel(level)
	}
	if err != nil {
		return logger, err
	}
	logger.SetLevel(logLevel)
	logger.SetOutput(io.MultiWriter(writers...))
	return logger, nil
}

func getLogPath(path string, name string) string {
	if strings.HasSuffix(path, string(os.PathSeparator)) {
		return path + name
	}
	return path + string(os.PathSeparator) + name
}