package logger

import (
	"log"
	"os"
)

type stdLogger struct {
	level  int
	logger *log.Logger
}

func NewStdLogger(level int) *stdLogger {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	return &stdLogger{
		level:  level,
		logger: logger,
	}
}

func (l *stdLogger) Debug(msg string, v ...any) {
	if LevelDebug < l.level {
		return
	}
	l.logger.Printf("[DEBUG] "+msg, v...)
}

func (l *stdLogger) Info(msg string, v ...any) {
	if LevelInfo < l.level {
		return
	}
	l.logger.Printf("[INFO]  "+msg, v...)
}

func (l *stdLogger) Warn(msg string, v ...any) {
	if LevelWarn < l.level {
		return
	}
	l.logger.Printf("[WARN]  "+msg, v...)
}

func (l *stdLogger) Error(msg string, v ...any) {
	if LevelError < l.level {
		return
	}
	l.logger.Printf("[ERROR] "+msg, v...)
}

func (l *stdLogger) Fatal(msg string, v ...any) {
	if LevelFatal < l.level {
		return
	}
	l.logger.Printf("[FATAL] "+msg, v...)
	os.Exit(1)
}
