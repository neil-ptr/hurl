package src

import (
	"bufio"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
)

const (
	METHOD   = 0
	URL      = 1
	PROTOCOL = 2
	NAME     = 0
	VALUE    = 1
)

func ParseHurlFile(filepath string, options Options) (HurlRequest, error) {
	f, err := os.OpenFile(filepath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	h := HurlRequest{}
	sc := bufio.NewScanner(f)

	// request line
	sc.Scan()
	requestLine := strings.TrimSpace(sc.Text())
	requestLineComponents := strings.Split(requestLine, " ")
	if len(requestLineComponents) > 2 {
		return HurlRequest{}, errors.New("Too many request line components")
	}
	if len(requestLineComponents) < 2 {
		return HurlRequest{}, errors.New("Not enough request line components")
	}

	parsedUrl, err := url.Parse(requestLineComponents[URL])
	h.URL = *parsedUrl
	h.Method = requestLineComponents[METHOD]

	// headers
	headerMap := make(map[string]string)
	for sc.Scan() && sc.Text() != "" {
		headerComponents := strings.Split(sc.Text(), ": ")
		headerName := headerComponents[NAME]
		headerVal := headerComponents[VALUE]

		headerMap[headerName] = headerVal
	}

	h.Headers = headerMap

	hostHeaderVal, exists := h.Headers["Host"]
	if exists && h.URL.Hostname() != hostHeaderVal {
		fmt.Println("Host header value does not match host in URL, using host in URL")
		h.Headers["Host"] = h.URL.Hostname()
	}

	h.Headers["User-Agent"] = "hurl/0.1.0"

	// body
	for sc.Scan() && len(sc.Bytes()) > 0 {
		h.Body = append(h.Body, sc.Bytes()...)
	}

	return h, nil
}
