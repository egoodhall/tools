package envars

import (
	"fmt"
	"sort"
	"strings"
)

func Println(prefix string, vars map[string]string) {
	for k, v := range vars {
		fmt.Println(render(prefix, k, v))
	}
}

func Render(prefix string, vars map[string]string) string {
	lines := make([]string, len(vars))
	var i int
	for k, v := range vars {
		lines[i] = render(prefix, k, v)
		i++
	}

	sort.Strings(lines)
	return strings.Join(lines, "\n")
}

func render(prefix, k, v string) string {
	sb := new(strings.Builder)
	if prefix != "" {
		sb.WriteString(prefix)
		sb.WriteRune(' ')
	}
	sb.WriteString(k)
	if v != "" {
		sb.WriteRune('=')
		sb.WriteString(v)
	}
	return sb.String()
}
