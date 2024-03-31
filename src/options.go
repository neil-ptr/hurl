package src

import (
	"flag"
)

type Config struct {
	Verbose        *bool
	BodyOutputPath *string
}

func InitConfig() Config {
	verbose := flag.Bool("v", false, "verbose output")
	bodyOutputPath := flag.String("b", "", "path to a file to output the response body")

	flag.Parse()

	return Config{
		Verbose:        verbose,
		BodyOutputPath: bodyOutputPath,
	}
}
