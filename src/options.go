package src

import (
	"flag"
)

type Options struct {
	Verbose        *bool
	BodyOutputPath *string
}

func InitOptions() Options {
	verbose := flag.Bool("v", false, "verbose output")
	bodyOutputPath := flag.String("b", "", "path to a file to output the response body")

	flag.Parse()

	return Options{
		Verbose:        verbose,
		BodyOutputPath: bodyOutputPath,
	}
}
