package src

import (
	"flag"
)

type Options struct {
	Verbose *bool
}

func InitOptions() Options {
	verbose := flag.Bool("v", false, "verbose output")

	flag.Parse()

	return Options{
		Verbose: verbose,
	}
}
