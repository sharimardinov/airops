package logger

import (
	"encoding/json"
	"io"
	"os"
	"time"
)

type JSONLogger struct {
	out io.Writer
}

func NewJSONLogger() *JSONLogger {
	return &JSONLogger{out: os.Stdout}
}

type LogEvent struct {
	Ts    string `json:"ts"`
	Level string `json:"level"`
	Msg   string `json:"msg"`

	RID    string `json:"rid,omitempty"`
	Method string `json:"method,omitempty"`
	Path   string `json:"path,omitempty"`

	Status int   `json:"status,omitempty"`
	Bytes  int   `json:"bytes,omitempty"`
	DurMS  int64 `json:"dur_ms,omitempty"`

	Route string `json:"route,omitempty"`
	IP    string `json:"ip,omitempty"`
	UA    string `json:"ua,omitempty"`
	Err   string `json:"err,omitempty"`
}

func (l *JSONLogger) Info(e LogEvent) {
	e.Level = "info"
	if e.Ts == "" {
		e.Ts = time.Now().UTC().Format(time.RFC3339Nano)
	}
	_ = json.NewEncoder(l.out).Encode(e)
}

func (l *JSONLogger) Error(e LogEvent) {
	e.Level = "error"
	if e.Ts == "" {
		e.Ts = time.Now().UTC().Format(time.RFC3339Nano)
	}
	_ = json.NewEncoder(l.out).Encode(e)
}
