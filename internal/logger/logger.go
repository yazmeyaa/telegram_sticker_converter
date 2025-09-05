package logger

import (
	"io"
	"log/slog"
	"os"
)

func NewLogger(w io.Writer, cmd string) *slog.Logger {
	log := slog.New(slog.NewJSONHandler(w, nil)).
		With(slog.Group(
			"program_info",
			slog.Int("pid", os.Getpid()),
		))

	return log
}
