package cli

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/egoodhall/tools/pkg/envars"
)

type SetCmd struct {
	Files envars.Files      `embed:""`
	Env   string            `arg:"" name:"env" required:""`
	Vars  map[string]string `arg:"" name:"vars"`
}

func (cmd *SetCmd) Run() error {
	data, err := os.ReadFile(cmd.Files.EnvFile(cmd.Env))
	if os.IsNotExist(err) {
		data = make([]byte, 0)
	} else if err != nil {
		return fmt.Errorf("load env file: %w", err)
	}

	vars := make(map[string]string)
	for _, line := range strings.Split(string(data), "\n") {
		if line != "" {
			name, value, found := strings.Cut(strings.Trim(line, " \n\t"), "=")
			if !found {
				return fmt.Errorf("invalid env line: %s", line)
			}
			vars[name] = value
		}
	}

	for k, v := range cmd.Vars {
		if _, ok := vars[k]; ok && v != "" {
			fmt.Fprintf(os.Stderr, "Overwriting %s\n", k)
		} else if _, ok := vars[k]; ok && v == "" {
			fmt.Fprintf(os.Stderr, "Clearing %s\n", k)
		} else if v != "" {
			fmt.Fprintf(os.Stderr, "Setting %s\n", k)
		}
		if v == "" {
			delete(vars, k)
		} else {
			vars[k] = v
		}
	}

	lines := make([]string, len(vars))
	var i int
	for k, v := range vars {
		lines[i] = fmt.Sprintf("%s=%s", k, v)
		i++
	}

	sort.Strings(lines)
	return os.WriteFile(cmd.Files.EnvFile(cmd.Env), []byte(strings.Join(lines, "\n")), 0600)
}
