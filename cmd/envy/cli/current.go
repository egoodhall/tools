package cli

import (
	"fmt"
	"strings"

	"github.com/egoodhall/tools/pkg/envars"
)

type CurrentCmd struct {
	Files envars.Files `embed:""`
	Envs  []string     `arg:"" name:"envs" optional:""`
}

func (cmd *CurrentCmd) Run() error {
	envs, err := cmd.Files.SelectedEnvs()
	if err != nil {
		return fmt.Errorf("load curent envs: %w", err)
	}

	fmt.Println(strings.Join(envs, "\n"))
	return nil
}
