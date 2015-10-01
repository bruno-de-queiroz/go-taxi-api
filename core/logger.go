package core

import (
	"code.google.com/p/log4go"
)

type Logger struct {
	*log4go.Logger
}

func NewLogger(config *LogConfig) *Logger {

	lvl := log4go.DEBUG

	switch config.Level {
	case "finest":
		lvl = log4go.FINEST
	case "fine":
		lvl = log4go.FINE
	case "trace":
		lvl = log4go.TRACE
	case "info":
		lvl = log4go.INFO
	case "warning":
		lvl = log4go.WARNING
	case "error":
		lvl = log4go.ERROR
	case "critical":
		lvl = log4go.CRITICAL
	}

	log := make(log4go.Logger)

	log.AddFilter("stdout", lvl, log4go.NewConsoleLogWriter())

	if config.File != "" {
		file := log4go.NewFileLogWriter(config.File, config.Rotate)

		if config.Format != "" {
			file.SetFormat(config.Format)
		}

		log.AddFilter("file", lvl, file)
	}

	return &Logger{&log}
}
