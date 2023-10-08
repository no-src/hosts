package main

import (
	"os"

	"github.com/no-src/hosts"
)

func main() {
	var search []string
	if len(os.Args) > 1 {
		search = os.Args[1:]
	}
	hosts.PrintHosts(search...)
}
