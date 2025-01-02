package logger

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

// Helper function to create a StdLogger with an in-memory buffer
func createLoggerWithBuffer(level int) (*stdLogger, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	logger := log.New(buf, "", log.Ldate|log.Ltime)
	return &stdLogger{level: level, logger: logger}, buf
}

func TestDebug(t *testing.T) {
	l, buf := createLoggerWithBuffer(LevelDebug)
	l.Debug("This is a debug message")
	output := buf.String()

	if !strings.Contains(output, "[DEBUG] This is a debug message") {
		t.Errorf("Expected log to contain '[DEBUG] This is a debug message', got: %s", output)
	}
}

func TestInfo(t *testing.T) {
	l, buf := createLoggerWithBuffer(LevelInfo)
	l.Info("This is an info message")
	output := buf.String()

	if !strings.Contains(output, "[INFO] This is an info message") {
		t.Errorf("Expected log to contain '[INFO] This is an info message', got: %s", output)
	}
}

func TestWarn(t *testing.T) {
	l, buf := createLoggerWithBuffer(LevelWarn)
	l.Warn("This is a warning message")
	output := buf.String()

	if !strings.Contains(output, "[WARN] This is a warning message") {
		t.Errorf("Expected log to contain '[WARN] This is a warning message', got: %s", output)
	}
}

func TestError(t *testing.T) {
	l, buf := createLoggerWithBuffer(LevelError)
	l.Error("This is an error message")
	output := buf.String()

	if !strings.Contains(output, "[ERROR] This is an error message") {
		t.Errorf("Expected log to contain '[ERROR] This is an error message', got: %s", output)
	}
}

func TestLevels(t *testing.T) {
	l, buf := createLoggerWithBuffer(LevelWarn)
	l.Debug("This is an debug message")
	output := buf.String()
	if output != "" {
		t.Errorf("Expected log to contain '', got: %s", output)
	}
	buf.Reset()

	l.Info("This is an info message")
	output = buf.String()
	if output != "" {
		t.Errorf("Expected log to contain '', got: %s", output)
	}
	buf.Reset()

	l.Warn("This is an warn message")
	output = buf.String()
	if !strings.Contains(output, "[WARN] This is an warn message") {
		t.Errorf("Expected log to contain '', got: %s", output)
	}
	buf.Reset()

	l.Error("This is an error message")
	output = buf.String()
	if !strings.Contains(output, "[ERROR] This is an error message") {
		t.Errorf("Expected log to contain '[ERROR] This is an error message', got: %s", output)
	}
}
