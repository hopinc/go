package leap

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// Logger is used to define the interface for a logger. It is un-opinionated, you will probably want to implement this
// interface for your logging client. By default, this client does not log.
type Logger interface {
	// Debug is used to log debug information.
	Debug(message string, metadata map[string]any)

	// Info is used to log information.
	Info(message string, metadata map[string]any)

	// Warn is used to log a warning.
	Warn(message string, metadata map[string]any)

	// Error is used to log an error. Note that the error object can be nil and this can just be error data.
	Error(message string, err error, metadata map[string]any)
}

// NopLogger is used to define a logger that does nothing.
type NopLogger struct{}

// Debug implements Logger.
func (NopLogger) Debug(string, map[string]any) {}

// Info implements Logger.
func (NopLogger) Info(string, map[string]any) {}

// Warn implements Logger.
func (NopLogger) Warn(string, map[string]any) {}

// Error implements Logger.
func (NopLogger) Error(string, error, map[string]any) {}

var _ Logger = NopLogger{}

// fmtLogger is used to define a logger that uses fmt. The logging level is controlled by the HOP_LOGGING_LEVEL environment
// variable. The default is "info". The levels are "debug", "info", "warn", and "error".
type fmtLogger struct {
	loggingLevel string
}

func formatMetadata(m map[string]any) string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	s := ""
	for _, k := range keys {
		v := m[k]
		switch x := v.(type) {
		case json.RawMessage:
			v = string(x)
		}
		s += " " + k + "=" + fmt.Sprint(v)
	}
	return s
}

// Debug implements Logger.
func (f fmtLogger) Debug(message string, metadata map[string]any) {
	if f.loggingLevel != "debug" {
		return
	}
	s := "[HOP] [DEBUG] " + message + formatMetadata(metadata)
	fmt.Println(s)
}

// Info implements Logger.
func (f fmtLogger) Info(message string, metadata map[string]any) {
	if f.loggingLevel != "debug" && f.loggingLevel != "info" {
		return
	}
	s := "[HOP] [INFO] " + message + formatMetadata(metadata)
	fmt.Println(s)
}

// Warn implements Logger.
func (f fmtLogger) Warn(message string, metadata map[string]any) {
	if f.loggingLevel != "debug" && f.loggingLevel != "info" && f.loggingLevel != "warn" {
		return
	}
	s := "[HOP] [WARN] " + message + formatMetadata(metadata)
	fmt.Println(s)
}

// Error implements Logger.
func (f fmtLogger) Error(message string, err error, metadata map[string]any) {
	if f.loggingLevel != "debug" && f.loggingLevel != "info" && f.loggingLevel != "warn" && f.loggingLevel != "error" {
		return
	}
	s := "[HOP] [ERROR] " + message
	if err != nil {
		s += " error=" + err.Error() + " stack=" + fmt.Sprintf("%+v", err)
	}
	s += formatMetadata(metadata)
	fmt.Println(s)
}

// NewFmtLogger creates a new fmt logger.
func NewFmtLogger() Logger {
	loggingInfo := os.Getenv("HOP_LOGGING_LEVEL")
	if loggingInfo != "debug" && loggingInfo != "info" && loggingInfo != "warn" && loggingInfo != "error" {
		loggingInfo = "info"
	}
	return fmtLogger{loggingInfo}
}
