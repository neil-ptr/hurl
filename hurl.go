package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/neil-and-void/hurl/src"
)

func main() {
	options := src.InitOptions()

	if len(os.Args) < 2 {
		fmt.Println("hurl: no hurl file provided")
		os.Exit(1)
	}

	hurlFilePath := os.Args[len(os.Args)-1]

	_, err := os.Stat(hurlFilePath)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Printf("hurl: file does not exist: %s\n", hurlFilePath)
		os.Exit(1)
	}

	f, err := os.OpenFile(hurlFilePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}

	hurlFile, err := src.ParseHurlFile(f)
	if err != nil {
		fmt.Println("hurl: ", err)
		os.Exit(1)
	}
	f.Close()

	if *options.Verbose {
		hurlFile.Output()
	}

	fmt.Printf("%+v\n", hurlFile)
}
