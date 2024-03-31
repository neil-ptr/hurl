package src

import (
	"bytes"
	"fmt"
	"net/http"
)

type HurlOutput struct {
	Config Config
}

func (h HurlOutput) OutputRequest(hurlFile HurlFile, req http.Request) error {
	buffer := bytes.Buffer{}

	statusLine := FormatStatusline(req)
	buffer.Write([]byte(statusLine))

	headers := FormatHeaders(req)
	buffer.Write(headers)

	// separate body with newline
	buffer.Write([]byte("\n"))

	if len(hurlFile.FilePaths) > 0 {
		buffer.Write(FormatFilePaths(hurlFile.FilePaths))
	} else {
		contentType := req.Header.Get("Content-Type")
		body, err := FormatBody(hurlFile.Body, contentType)
		if err != nil {
			return err
		}

		buffer.Write(body)
	}

	fmt.Printf("%s\n", buffer.String())

	return nil
}

func (h HurlOutput) OutputResponse(res http.Response) {}
