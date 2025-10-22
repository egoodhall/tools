package main

import (
	"errors"
	"flag"
	"log/slog"
	"os"
)

func parseFlags() (Flags, error) {
	var (
		logLevel LogLevel
		write    bool
		prune    bool
		keysFile string
	)

	flag.Var(logLevel, "log", "Log level")
	flag.BoolVar(&write, "write", true, "Whether the synced values should be written")
	flag.BoolVar(&prune, "prune", true, "Whether to prune keys that are no longer found")
	flag.StringVar(&keysFile, "keys", "$HOME/.ssh/authorized_keys", "Path to the keys file")
	flag.Parse()

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel.Get(),
	})))

	keysFile = os.ExpandEnv(keysFile)

	if len(flag.Args()) == 0 {
		return Flags{}, errors.New("no users specified")
	}

	return Flags{
		Write:    write,
		Prune:    prune,
		KeysFile: keysFile,
	}, nil
}

type Flags struct {
	Write    bool
	Prune    bool
	KeysFile string
}

var _ flag.Value = new(LogLevel)

type LogLevel struct {
	level slog.Level
}

func (l LogLevel) String() string {
	return l.level.String()
}

func (l LogLevel) Set(v string) error {
	return l.level.UnmarshalText([]byte(v))
}

func (l LogLevel) Get() slog.Level {
	return l.level
}
