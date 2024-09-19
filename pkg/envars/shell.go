package envars

import (
	"os"
	"strings"
)

func LoadEnv() (map[string]string, error) {
	vars := make(map[string]string)

	if current, ok := os.LookupEnv("EMV_MANAGED_VARS"); ok {
		names := strings.Split(current, ":")
		for _, name := range names {
			if val, ok := os.LookupEnv(name); ok {
				vars[name] = val
			}
		}
	}

	return vars, nil
}

func IsSame(current, next map[string]string) bool {
	for k, v := range current {
		if nval, ok := next[k]; !ok || v != nval {
			return false
		}
	}
	for k, v := range next {
		if cval, ok := current[k]; !ok || v != cval {
			return false
		}
	}
	return true
}
