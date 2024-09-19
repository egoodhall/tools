package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/egoodhall/tools/pkg/envars"
)

type UnsetCmd struct {
	Files envars.Files `embed:""`
	Env   string       `arg:"" name:"env" required:""`
	All   bool         `name:"all" short:"a" xor:"all"`
	Vars  []string     `arg:"" name:"vars" optional:"" xor:"all"`
}

func (cmd *UnsetCmd) Run() error {
	if cmd.All {
		return os.WriteFile(cmd.Files.EnvFile(cmd.Env), []byte{}, 0600)
	}

	data, err := os.ReadFile(cmd.Files.EnvFile(cmd.Env))
	if err != nil {
		return fmt.Errorf("load env file: %w", err)
	}

	vars := make(map[string]string)
	for _, line := range strings.Split(string(data), "\n") {
		name, value, found := strings.Cut(line, "=")
		if !found {
			return fmt.Errorf("invalid env line: %s", line)
		}
		vars[name] = value
	}

	for _, v := range cmd.Vars {
		if _, ok := vars[v]; ok {
			fmt.Fprintf(os.Stderr, "Removing %s\n", v)
		}
		delete(vars, v)
	}

	return os.WriteFile(cmd.Files.EnvFile(cmd.Env), []byte(envars.Render("", vars)), 0600)
}
