package main

import (
	"github.com/kartverket/skyline/cmd"
	"log/slog"
	"os"
)

func main() {
	// Configure logging
	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(l)

	// Launch
	cmd.Execute()
}
