package api

import (
	"fmt"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func NewRequestLogger(logger *log.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&RequestLogger{logger})
}

type RequestLogger struct {
	Logger *log.Logger
}

func (l *RequestLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &StructuredLoggerEntry{Logger: log.NewEntry(l.Logger)}
	logFields := log.Fields{}
	logFields["uri"] = fmt.Sprintf("%s%s", r.Host, r.RequestURI)
	entry.Logger = entry.Logger.WithFields(logFields)
	entry.Logger.Infoln("request started")
	return entry
}

type StructuredLoggerEntry struct {
	Logger log.FieldLogger
}

func (l *StructuredLoggerEntry) Write(status, bytes int, elapsed time.Duration) {
	l.Logger = l.Logger.WithFields(log.Fields{
		"status":     status,
		"lengths":    bytes,
		"elapsed_ms": float64(elapsed.Nanoseconds()) / float64(time.Millisecond),
	})

	l.Logger.Infoln("request complete")
}

func (l *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.Logger = l.Logger.WithFields(log.Fields{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	})
}
