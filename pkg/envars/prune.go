package envars

import (
	"fmt"
	"os"
	"strings"

	"log/slog"
)

func PruneUnusedEnvs(files Files) error {
	envs, err := files.AllEnvs()
	if err != nil {
		return fmt.Errorf("load all envs: %w", err)
	}

	for env := range envs {
		vars, err := files.Load(env)
		if err != nil {
			return fmt.Errorf("load env vars: %w", err)
		}

		if len(vars) != 0 {
			delete(envs, env)
		}
	}

	activeEnvs, err := files.SelectedEnvs()
	if err != nil {
		return err
	}

	prunedActiveEnvs := make([]string, 0)

	for pruneEnv := range envs {
		if err := os.Remove(files.EnvFile(pruneEnv)); err != nil {
			slog.Error("couldn't prune env file", "env", pruneEnv, "error", err)
		}
	}

	for _, activeEnv := range activeEnvs {
		if _, ok := envs[activeEnv]; !ok {
			prunedActiveEnvs = append(prunedActiveEnvs, activeEnv)
		}
	}

	if err := os.WriteFile(files.CurrentEnvsFile(), []byte(strings.Join(prunedActiveEnvs, " ")), 0600); err != nil {
		return fmt.Errorf("update selected envs file: %w", err)
	}
	return nil
}
