package cli

import (
	"fmt"

	"github.com/egoodhall/tools/pkg/envars"
)

type ShowCmd struct {
	Files envars.Files `embed:""`
	Envs  []string     `arg:"" name:"envs" optional:""`
}

func (cmd *ShowCmd) Run() error {
	var env map[string]string
	var err error
	if len(cmd.Envs) == 0 {
		env, err = envars.LoadEnv()
	} else {
		env, err = cmd.Files.Load(cmd.Envs...)
	}

	if err != nil {
		return fmt.Errorf("load envs: %w", err)
	}

	fmt.Println(envars.Render("", env))
	return nil
}
