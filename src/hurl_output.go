package src

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"time"
)

var LOADING_CHARS = [...]string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

type HurlOutput struct {
	Config HurlConfig
}

func WaitForHttpRequest(req *http.Request) (*http.Response, error) {
	errCh := make(chan error)
	resCh := make(chan *http.Response)

	go func() {
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			errCh <- err
			return
		}

		resCh <- res
	}()

	i := 0
	for {
		select {
		case res := <-resCh:
			ClearSpinner()
			return res, nil
		case err := <-errCh:
			ClearSpinner()
			return nil, err
		default:
			PrintSpinner(i)
			time.Sleep(50 * time.Millisecond)
		}
		i = (i + 1) % len(LOADING_CHARS)
	}
}

func PrintSpinner(i int) {
	fmt.Printf("=== sending %s ===\r", LOADING_CHARS[i])
}

func ClearSpinner() {
	fmt.Print("\r")
}

func (h HurlOutput) OutputRequest(hurlFile HurlFile, req http.Request) error {
	buffer := bytes.Buffer{}

	requestLine := FormatRequestLine(req)
	buffer.Write([]byte(requestLine))

	// add this in manually since not set in http.Request for some reason
	req.Header.Set("Host", req.Host)

	headers := FormatHeaders(req.Header, ">")
	buffer.Write(headers)

	// separate body with newline
	if len(hurlFile.Body) == 0 && len(hurlFile.FileEmbed) == 0 && len(hurlFile.MultipartFormData) == 0 {
		fmt.Printf("%s\n", buffer.String())
		return nil
	}

	buffer.Write([]byte("\n"))

	contentType := req.Header.Get("Content-Type")

	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil
	}

	if mediaType == "multipart/form-data" {
		body, err := FormatMultiPart(hurlFile.MultipartFormData, hurlFile.MultipartBoundary)
		if err != nil {
			return err
		}

		buffer.Write(body)

	} else if hurlFile.FileEmbed != "" {
		fileEmbed := fmt.Sprintf("%s\n", FormatFileEmbed(hurlFile.FileEmbed))
		buffer.Write([]byte(fileEmbed))

	} else {
		body, err := FormatBody(hurlFile.Body, mediaType)
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

	headers := FormatHeaders(res.Header, "<")
	buffer.Write([]byte(headers))

	// separate body with newline
	buffer.Write([]byte("\n"))

	contentType := res.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	bodyOutputPath := h.Config.BodyOutputPath
	if len(bodyOutputPath) > 0 {
		if mediaType == "application/json" {
			prettified, err := PrettifyJson(bodyBytes)
			if err != nil {
				return err
			}

			bodyBytes = prettified
		}

		err := os.WriteFile(bodyOutputPath, bodyBytes, 0644)
		if err != nil {
			return err
		}

		title := FormatFilePathsTitle()
		buffer.Write([]byte(title))

		filePaths := FormatFileEmbed(bodyOutputPath)
		buffer.Write(filePaths)

		fmt.Printf("%s\n", buffer.String())

		return nil
	}

	body, err := FormatBody(bodyBytes, mediaType)
	if err != nil {
		return err
	}

	buffer.Write(body)

	fmt.Printf("%s\n", buffer.String())

	return nil
}
