package src

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
)

type HurlRequest struct {
	Method  string
	URL     url.URL
	Headers map[string]string
	Body    []byte
	// FilePaths []string
}

func (h HurlRequest) HttpRequest() (*http.Request, error) {
	req, err := http.NewRequest(h.Method, h.URL.String(), nil)
	if err != nil {
		return nil, err
	}

	for name, value := range h.Headers {
		req.Header.Add(name, value)
	}

	// add body
	if len(h.Body) > 0 {
		req.Body = io.NopCloser(bytes.NewReader(h.Body))
		req.ContentLength = int64(len(h.Body))
	}

	return req, nil
}
