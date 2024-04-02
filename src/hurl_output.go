package src

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

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
		i += 1
	}
}

func PrintSpinner(i int) {
	loadingChar := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

	i = i % len(loadingChar)
	fmt.Printf("=== loading %s ===\r", loadingChar[i])
}

func ClearSpinner() {
	fmt.Print("\r")
}

func (h HurlOutput) OutputRequest(hurlFile HurlFile, req http.Request) error {
	buffer := bytes.Buffer{}

	requestLine := FormatRequestLine(req)
	buffer.Write([]byte(requestLine))

	headers := FormatHeaders(req.Header, ">")
	buffer.Write(headers)

	// separate body with newline
	if len(hurlFile.Body) == 0 {
		fmt.Printf("%s\n", buffer.String())
		return nil
	}

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

	headers := FormatHeaders(res.Header, "<")
	buffer.Write([]byte(headers))

	// separate body with newline
	buffer.Write([]byte("\n"))

	contentType := res.Header.Get("Content-Type")
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	bodyOutputPath := h.Config.BodyOutputPath
	if len(bodyOutputPath) > 0 {
		err := os.WriteFile(bodyOutputPath, bodyBytes, 0644)
		if err != nil {
			return err
		}

		title := FormatFilePathsTitle()
		buffer.Write([]byte(title))

		filePaths := FormatFilePaths([]string{bodyOutputPath})
		buffer.Write(filePaths)

		fmt.Printf("%s\n", buffer.String())

		return nil
	}

	body, err := FormatBody(bodyBytes, contentType)
	if err != nil {
		return err
	}

	buffer.Write(body)

	fmt.Printf("%s\n", buffer.String())

	return nil
}
