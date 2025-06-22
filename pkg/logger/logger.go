package logger

import (
	"log"
	"os"

	kitlog "github.com/go-kit/log"
)

// NewLogger creates a new logger.
func NewLogger() kitlog.Logger {
	logger := kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stdout))

	log.SetOutput(kitlog.NewStdlibAdapter(logger))

	return kitlog.WithPrefix(logger, "ts", kitlog.DefaultTimestamp, "file", kitlog.DefaultCaller)
}
