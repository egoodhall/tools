package main

import (
	"log/slog"

	"github.com/alecthomas/kong"
	"github.com/egoodhall/tools/cmd/sshkeyd/cli"
	"github.com/egoodhall/tools/pkg/logging"
)

func main() {
	ktx := kong.Parse(new(Cli), kong.DefaultEnvars("SSHKEYD"))
	ktx.FatalIfErrorf(ktx.Run())
}

type Cli struct {
	Level slog.Level `name:"level" short:"l" required:"" default:"info"`

	Sync cli.SyncCmd `name:"sync" cmd:""`
}

func (cli *Cli) AfterApply() error {
	slog.SetDefault(slog.New(logging.NewHandler("", cli.Level)))
	return nil
}
