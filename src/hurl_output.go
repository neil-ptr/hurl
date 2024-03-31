package src

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type HurlOutput struct {
	Config Config
}

func (h HurlOutput) OutputRequest(hurlFile HurlFile, req http.Request) error {
	buffer := bytes.Buffer{}

	requestLine := FormatRequestLine(req)
	buffer.Write([]byte(requestLine))

	headers := FormatHeaders(req.Header)
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

func (h HurlOutput) OutputResponse(res http.Response) error {
	buffer := bytes.Buffer{}

	statusLine := FormatStatusLine(res)
	buffer.Write([]byte(statusLine))

	headers := FormatHeaders(res.Header)
	buffer.Write([]byte(headers))

	// separate body with newline
	buffer.Write([]byte("\n"))

	if h.Config.BodyOutputPath == nil {

		// output to file at path
		return nil
	}

	contentType := res.Header.Get("Content-Type")
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	body, err := FormatBody(bodyBytes, contentType)
	if err != nil {
		return err
	}

	buffer.Write(body)

	fmt.Printf("%s\n", buffer.String())

	return nil
}
