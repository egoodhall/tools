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
		help    bool
	)

	flag.BoolVar(&help, "h", false, "Show this help message")
	flag.Parse()
	comment = strings.Join(flag.Args(), " ")

	if help || comment == "" {
		flag.Usage()
		return
	}

	comment = strings.TrimSpace(comment)
	pipes := strings.Repeat("-", len(comment))

	fmt.Printf(strings.TrimSpace(template)+"\n", pipes, comment, pipes)
}
