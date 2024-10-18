package main

import (
	"log/slog"
	"os"

	"github.com/alecthomas/kong"
)

var programLevel = new(slog.LevelVar)

type cmdContext struct {
	Logger *slog.Logger
}

type kongApp struct {
	LogLevel int `short:"l" help:"Log level: 0 (debug), 1 (info), 2 (warn), 3 (error)" default:"0"`

	Server serverCmd `cmd:"" help:"Start the app server."`
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel}))
	slog.SetDefault(logger)

	app := &kongApp{}
	kongCtx := kong.Parse(app,
		kong.Description("A simple application."),
		kong.UsageOnError(),
		kong.Vars{
			"version": "0.0.1",
		},
	)
	switch app.LogLevel {
	case 0:
		programLevel.Set(slog.LevelDebug)
	case 1:
		programLevel.Set(slog.LevelInfo)
	case 2:
		programLevel.Set(slog.LevelWarn)
	case 3:
		programLevel.Set(slog.LevelError)
	}
	err := kongCtx.Run(&cmdContext{Logger: logger})
	kongCtx.FatalIfErrorf(err)
}
