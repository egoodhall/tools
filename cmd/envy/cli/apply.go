package cli

import (
	"fmt"
	"os"
	"strings"

	"log/slog"

	"github.com/egoodhall/tools/pkg/envars"
)

type ApplyCmd struct {
	Files envars.Files `embed:""`
}

func (cmd *ApplyCmd) Run() error {
	if err := envars.PruneUnusedEnvs(cmd.Files); err != nil {
		return fmt.Errorf("prune unused envs: %w", err)
	}

	currentVars, err := envars.LoadEnv()
	if err != nil {
		return fmt.Errorf("load applied env vars: %w", err)
	}

	vars, err := cmd.Files.Load()
	if err != nil {
		return fmt.Errorf("load file env vars: %w", err)
	}

	slog.Debug("Compare applied & expected env variables", "applied", currentVars, "file", vars)
	if envars.IsSame(currentVars, vars) {
		slog.Info("Applied and expected environments match. No changes needed")
		return nil
	}
	slog.Info("Expected environment changed. Updating.")

	if current, ok := os.LookupEnv("EMV_MANAGED_VARS"); ok {
		names := strings.Split(current, ":")
		unset := make(map[string]string)
		for _, name := range names {
			unset[name] = ""
		}
		envars.Println("unset", unset)
		fmt.Println("unset EMV_MANAGED_VARS")
	}

	var i int
	names := make([]string, len(vars))
	for name := range vars {
		names[i] = name
		i++
	}

	if len(names) == 0 {
		return nil
	}

	vars["EMV_MANAGED_VARS"] = strings.Join(names, ":")
	envars.Println("export", vars)
	return nil
}
