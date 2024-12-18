package main

import (
	"flag"
	"fmt"
	"strings"
)

const template = `
/*-%s-*
 | %s |
 *-%s-*/
`

func main() {
	var (
		comment string
		indent  int
		help    bool
	)

	flag.BoolVar(&help, "h", false, "Show this help message")
	flag.IntVar(&indent, "i", 2, "Indentation level")
	flag.Parse()
	comment = strings.Join(flag.Args(), " ")

	if help || comment == "" {
		flag.Usage()
		return
	}

	comment = strings.TrimSpace(comment)
	pipes := strings.Repeat("-", len(comment))
	tmpl := strings.TrimLeft(template, "\n")

	fmt.Printf(indentString(tmpl, indent), pipes, comment, pipes)
}

func indentString(str string, n int) string {
	lines := strings.Split(str, "\n")
	for i, line := range lines {
		lines[i] = strings.Repeat(" ", n) + line
	}
	return strings.Join(lines, "\n")
}
