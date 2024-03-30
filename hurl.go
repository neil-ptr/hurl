package main

import (
	"errors"
	"fmt"
	"net/http"
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

	hurlRequest, err := src.ParseHurlFile(hurlFilePath, options)
	if err != nil {
		fmt.Println("hurl: ", err)
		os.Exit(1)
	}

	req, err := hurlRequest.HttpRequest()
	if err != nil {
		fmt.Println("hurl: ", err)
		os.Exit(1)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("hurl: ", err)
		os.Exit(1)
	}

	formattedReponse, err := src.FormatResponse(res)
	if err != nil {
		fmt.Println("hurl: ", err)
		os.Exit(1)
	}

	if *options.Verbose == true {
		formattedRequest, err := src.FormatRequest(req)
		if err != nil {
			fmt.Println("hurl: ", err)
			os.Exit(1)
		}

		fmt.Printf("%s", formattedRequest)
	}

	fmt.Printf("%s\n", formattedReponse)
}
