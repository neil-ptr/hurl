package main

import (
	"fmt"
	"os"

	"github.com/neil-and-void/hurl/src"
)

func main() {
	config, err := src.InitConfig()

	if config.Version {
		fmt.Println("v0.4.1")
		os.Exit(0)
	}

	if err != nil {
		fmt.Printf("hurl: %s\n", err.Error())
		os.Exit(1)
	}

	hurlOutput := src.HurlOutput{Config: config}

	if len(os.Args) < 2 {
		fmt.Println("hurl: no hurl file provided")
		os.Exit(1)
	}

	hurlFilePath := os.Args[len(os.Args)-1]

	_, err = os.Stat(hurlFilePath)
	if err != nil {
		fmt.Printf("hurl: file does not exist: %s\n", hurlFilePath)
		os.Exit(1)
	}

	f, err := os.OpenFile(hurlFilePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		fmt.Printf("hurl: %s\n", err.Error())
		os.Exit(1)
	}

	hurlFile, err := src.ParseHurlFile(f)
	if err != nil {
		fmt.Printf("hurl: %s\n", err.Error())
		os.Exit(1)
	}
	f.Close()

	req, err := hurlFile.NewRequest()
	if err != nil {
		fmt.Printf("hurl: %s\n", err.Error())
		os.Exit(1)
	}

	if config.Verbose {
		err = hurlOutput.OutputRequest(hurlFile, *req)
		if err != nil {
			fmt.Printf("hurl: %s\n", err.Error())
			os.Exit(1)
		}
	}

	res, err := src.WaitForHttpRequest(req)
	if err != nil {
		fmt.Printf("hurl: %s\n", err.Error())
		os.Exit(1)
	}

	err = hurlOutput.OutputResponse(*res)
	if err != nil {
		fmt.Printf("hurl: %s\n", err.Error())
		os.Exit(1)
	}
}
