package global

import (
	"log/slog"
	"os"
)

func SetDefaultLogger() {
	options := &slog.HandlerOptions{
		Level: ProgramLevel,
	}

	ProgramLevel.Set(slog.LevelInfo)
	Logger = slog.New(slog.NewTextHandler(os.Stderr, options))
	slog.SetDefault(Logger)
}
