package main

import (
	"os"

	"log/slog"

	"github.com/alecthomas/kong"
	"github.com/egoodhall/tools/cmd/envy/cli"
	"github.com/egoodhall/tools/pkg/logging"
	"github.com/willabides/kongplete"
)

var Cli struct {
	LogConfig   logging.Handler              `embed:"" prefix:"log-"`
	Set         cli.SetCmd                   `cmd:"" name:"set" help:"set environment variables"`
	Unset       cli.UnsetCmd                 `cmd:"" name:"unset" help:"unset environment variables"`
	Show        cli.ShowCmd                  `cmd:"" name:"show" help:"show the variables in an environment"`
	Select      cli.SelectCmd                `cmd:"" name:"select" help:"select the environments to use"`
	Current     cli.CurrentCmd               `cmd:"" name:"current" help:"show the currently selected environments"`
	Apply       cli.ApplyCmd                 `cmd:"" name:"apply" help:"output shell commands to apply the currently selected environment"`
	Completions kongplete.InstallCompletions `cmd:"" name:"install-completions" help:"install shell completions"`
}

func main() {
	k := kong.Must(&Cli,
		kong.DefaultEnvars("ENVY"),
	)
	kongplete.Complete(k)
	ktx, err := k.Parse(os.Args[1:])
	slog.SetDefault(slog.New(&Cli.LogConfig))
	k.FatalIfErrorf(err)
	k.FatalIfErrorf(ktx.Run())
}
