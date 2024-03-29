package parser

import (
	"bufio"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
)

var CRLF = []byte{13, 10} // [c]arriage [r]eturn, [l]ine [f]eed

const (
	METHOD   = 0
	URL      = 1
	PROTOCOL = 2
	NAME     = 0
	VALUE    = 1
)

type HurlRequest struct {
	Method  string
	URL     *url.URL
	Headers map[string]string
	Body    []byte
	// FilePaths []string
}

func (h HurlRequest) GetRawRequest() []byte {
	rawRequest := []byte{}

	requestLine := fmt.Sprintf("%s %s HTTP/1.1\r\n", h.Method, h.Path)
	rawRequest = append(rawRequest, []byte(requestLine)...)

	for key := range h.Headers {
		headerVal, _ := h.Headers[key]

		header := fmt.Sprintf("%s: %s\r\n", key, headerVal)
		rawRequest = append(rawRequest, []byte(header)...)
	}

	rawRequest = append(rawRequest, CRLF...)

	if len(h.Body) == 0 {
		return rawRequest
	}

	// TODO: process body

	return rawRequest
}

func ParseHurlFile(filepath string) (HurlRequest, error) {
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
	h.URL = parsedUrl
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

	host, ok := headerMap["Host"]
	if !ok {
		panic("`HOST` header is missing")
	}

	h.Host = host

	// TODO: body

	return h, nil
}
