package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/neil-and-void/hurl/src"
)

func main() {
	config, err := src.InitConfig()
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
		fmt.Printf("hurl: %s", err.Error())
		os.Exit(1)
	}

	hurlFile, err := src.ParseHurlFile(f)
	if err != nil {
		fmt.Printf("hurl: %s", err.Error())
		os.Exit(1)
	}
	f.Close()

	req, err := hurlFile.NewRequest()

	if config.Verbose {
		err = hurlOutput.OutputRequest(hurlFile, *req)
		if err != nil {
			fmt.Printf("hurl: %s", err.Error())
			os.Exit(1)
		}
	}

	res, err := http.DefaultClient.Do(req)

	err = hurlOutput.OutputResponse(*res)
	if err != nil {
		fmt.Printf("hurl: %s", err.Error())
		os.Exit(1)
	}

}
