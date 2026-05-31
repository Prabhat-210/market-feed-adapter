package logger

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

func NewLogger(svc, env, logLevel string) zerolog.Logger {
	lvl, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		lvl = zerolog.InfoLevel
	}
	writer := getWriter(env)
	return zerolog.New(writer).
		Level(lvl).
		With().
		Timestamp().
		Caller().
		Str("service", svc).
		Str("env", env).
		Logger()
}

func getWriter(env string) io.Writer {
	switch strings.ToLower(env) {

	case "local", "dev":
		return zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	default:
		return os.Stdout //no need to set TimeFormat In prod, logs are JSON and zerolog handles timestamp automatically.
	}
}
