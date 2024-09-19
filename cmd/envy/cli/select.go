package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/egoodhall/tools/pkg/envars"
)

type SelectCmd struct {
	Files envars.Files `embed:""`
	Clear bool         `name:"clear" xor:"clear"`
	Envs  []string     `arg:"" name:"env" optional:"" xor:"clear"`
}

func (cmd *SelectCmd) Run() error {
	if cmd.Clear {
		return os.WriteFile(cmd.Files.CurrentEnvsFile(), []byte{}, 0600)
	}

	if len(cmd.Envs) == 0 {
		envs, err := cmd.Files.SelectedEnvs()
		if err != nil {
			return err
		}
		fmt.Println(strings.Join(envs, " "))
		return nil
	}

	envs, err := cmd.Files.AllEnvs()
	if err != nil {
		return err
	}

	for _, env := range cmd.Envs {
		if _, ok := envs[env]; !ok {
			return fmt.Errorf("unknown env: %s", env)
		}
	}

	return os.WriteFile(cmd.Files.CurrentEnvsFile(), []byte(strings.Join(cmd.Envs, " ")), 0600)
}
