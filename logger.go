package main

import (
	log "github.com/sirupsen/logrus"
	"os"
)

type customLogger struct {
	appName   string
	appVer    string
	formatter log.JSONFormatter
}

func (l customLogger) Format(entry *log.Entry) ([]byte, error) {
	entry.Data["appName"] = l.appName
	entry.Data["appVer"] = l.appVer
	return l.formatter.Format(entry)
}

func GetLogger(ctx *Context) *log.Logger {
	var logger = log.New()
	logger.Out = os.Stdout
	logger.SetFormatter(customLogger{
		appName:   appName,
		appVer:    appVer,
		formatter: log.JSONFormatter{},
	})

	logger.Level = ctx.config.Logger.level

	return logger
}
