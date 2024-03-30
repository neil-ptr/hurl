package main

import (
	"log"
	"net/http"
	"os"

	"github.com/neil-and-void/hurl/src"
)

func main() {
	args := os.Args[1:]

	hurlRequest, err := src.ParseHurlFile(args[0])
	if err != nil {
		log.Fatal("hurl: ", err)
	}

	req, err := hurlRequest.HttpRequest()
	if err != nil {
		log.Fatal("hurl: ", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("hurl: ", err)
	}

	src.WriteResponse(res)
}
