package main

import (
	"os"
	"os/exec"

	"log/slog"

	"github.com/alecthomas/kong"
	"github.com/egoodhall/tools/pkg/logging"
)

func main() {
	ctx := kong.Parse(new(Cli), kong.DefaultEnvars("INDIR"))
	ctx.FatalIfErrorf(ctx.Run())
}

type Cli struct {
	LogLevel slog.Level `name:"log" short:"l" required:"" default:"INFO"`
	Dir      struct {
		Dir string `name:"directory" arg:"" required:""`
		Do  struct {
			Command string   `name:"command" arg:"" required:""`
			Args    []string `name:"args" arg:"" optional:""`
		} `cmd:""`
	} `arg:"" name:"directory"`
}

func (cli *Cli) AfterApply() error {
	slog.SetDefault(slog.New(logging.NewHandler("", cli.LogLevel)))
	return nil
}

func (cli *Cli) Run() error {
	slog.Debug("running command in directory", "directory", cli.Dir, "command", cli.Dir.Do.Command, "args", cli.Dir.Do.Args)
	if err := os.Chdir(cli.Dir.Dir); err != nil {
		return err
	}

	cmd := exec.Command(cli.Dir.Do.Command, cli.Dir.Do.Args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
