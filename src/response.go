package src

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/alecthomas/chroma/v2/quick"
	"github.com/fatih/color"
)

type HurlResponse struct {
	Response *http.Response
	Duration time.Duration
}

func (h HurlResponse) Format() ([]byte, error) {
	formattedResponse := []byte{}

	protocol := formatProtocol(h.Response.Proto)
	status := formatStatusCode(h.Response.StatusCode, h.Response.Status)
	requestLine := fmt.Sprintf("< %s%s\n", protocol, status)

	formattedResponse = append(formattedResponse, []byte(requestLine)...)

	coloredHeaderName := color.New(color.FgYellow).SprintFunc()
	for name, value := range h.Response.Header {
		formattedResponse = append(formattedResponse, []byte(fmt.Sprintf("< %s: %s\n", coloredHeaderName(name), strings.Join(value, "")))...)
	}

	// separate headers from body
	formattedResponse = append(formattedResponse, []byte("\n")...)
	body, err := io.ReadAll(h.Response.Body)
	if err != nil {
		return []byte{}, err
	}

	formatter := "" // raw text
	contentType := h.Response.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		formatter = "json"
	} else if strings.Contains(contentType, "text/html") {
		formatter = "html"
	}

	if len(body) == 0 {
		return []byte{}, nil
	}

	buffer := bytes.Buffer{}
	err = quick.Highlight(&buffer, string(body), formatter, "terminal", "")
	formattedResponse = append(formattedResponse, buffer.Bytes()...)

	return formattedResponse, nil
}

func (h HurlResponse) Output(w io.Writer) error {
	return nil
}
