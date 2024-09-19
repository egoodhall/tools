package envars

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Files struct {
	Root string `name:"root-dir" required:"" default:"$HOME/.env"`
}

func (f Files) Load(envs ...string) (map[string]string, error) {
	if len(envs) == 0 {
		selectedEnvs, err := f.SelectedEnvs()
		if err != nil {
			return nil, fmt.Errorf("selected envs: %w", err)
		}
		envs = selectedEnvs
	}

	env := make(map[string]string)
	for _, ef := range envs {
		if ef == "" {
			continue
		}

		data, err := os.ReadFile(f.EnvFile(ef))
		if err != nil {
			return nil, err
		}

		data = bytes.Trim(data, "\n\t ")
		if len(data) == 0 {
			continue
		}

		for _, line := range strings.Split(string(data), "\n") {
			name, value, found := strings.Cut(line, "=")
			if !found {
				return nil, fmt.Errorf("invalid env line: %s", line)
			}
			env[name] = value
		}
	}
	return env, nil
}

func (f Files) EnvFile(env string) string {
	return os.ExpandEnv(filepath.Join(f.Root, strings.Trim(env, " \n\t")))
}

func (f Files) CurrentEnvsFile() string {
	return f.EnvFile(".current")
}

func (f Files) SelectedEnvs() ([]string, error) {
	data, err := os.ReadFile(f.CurrentEnvsFile())
	if err != nil {
		return nil, fmt.Errorf("load current envs: %w", err)
	}
	return strings.Split(strings.Trim(string(data), "\n\t "), " "), nil
}

func (f Files) AllEnvs() (map[string]struct{}, error) {
	des, err := os.ReadDir(os.ExpandEnv(f.Root))
	if err != nil {
		return nil, fmt.Errorf("read envs dir: %w", err)
	}

	envs := make(map[string]struct{})
	for _, de := range des {
		if de.IsDir() || de.Name() == ".current" {
			continue
		}

		envs[de.Name()] = struct{}{}
	}

	return envs, nil
}
