package src

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/alecthomas/chroma/v2/quick"
	"github.com/fatih/color"
)

func FormatRequest(req *http.Request) ([]byte, error) {
	formattedRequest := []byte{}

	protocol := formatProtocol(req.Proto)
	method := formatMethod(req.Method)
	path := formatPath(req.URL.Path)

	requestLine := fmt.Sprintf("> %s%s%s\n", method, path, protocol)

	formattedRequest = append(formattedRequest, []byte(requestLine)...)

	for name, value := range req.Header {
		coloredHeaderName := color.New(color.FgYellow).SprintFunc()
		formattedRequest = append(formattedRequest, []byte(fmt.Sprintf("> %s: %s\n", coloredHeaderName(name), strings.Join(value, "")))...)
	}

	formattedRequest = append(formattedRequest, []byte("> \n")...)

	return formattedRequest, nil
}

func FormatResponse(res *http.Response) ([]byte, error) {
	formattedResponse := []byte{}

	protocol := formatProtocol(res.Proto)
	status := formatStatusCode(res.StatusCode, res.Status)
	requestLine := fmt.Sprintf("< %s%s\n", protocol, status)

	formattedResponse = append(formattedResponse, []byte(requestLine)...)

	coloredHeaderName := color.New(color.FgYellow).SprintFunc()
	for name, value := range res.Header {
		formattedResponse = append(formattedResponse, []byte(fmt.Sprintf("< %s: %s\n", coloredHeaderName(name), strings.Join(value, "")))...)
	}

	// separate headers from body
	formattedResponse = append(formattedResponse, []byte("\n")...)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}

	formatter := "" // raw text
	contentType := res.Header.Get("Content-Type")
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

func formatMethod(method string) string {
	formatted := fmt.Sprintf(" %s ", method)

	switch method {
	case "GET":
		return fmt.Sprintf("%s", color.New(color.BgGreen, color.FgBlack).Sprint(formatted))

	case "POST":
		return fmt.Sprintf("%s", color.New(color.BgBlue, color.FgBlack).Sprint(formatted))

	case "PUT":
		return fmt.Sprintf("%s", color.New(color.BgYellow, color.FgBlack).Sprint(formatted))

	case "PATCH":
		return fmt.Sprintf("%s", color.New(color.BgYellow, color.FgBlack).Sprint(formatted))

	case "DELETE":
		return fmt.Sprintf("%s", color.New(color.BgRed, color.FgBlack).Sprint(formatted))

	default:
		return fmt.Sprintf("%s", color.New(color.BgMagenta, color.FgBlack).Sprint(formatted))
	}
}

func formatPath(path string) string {
	coloredPath := color.New(color.BgWhite, color.FgBlack).SprintFunc()
	formatted := fmt.Sprintf(" %s ", path)
	return fmt.Sprintf("%s", coloredPath(formatted))
}

func formatStatusCode(statusCode int, status string) string {
	coloredStatus := color.New(color.BgWhite, color.FgBlack).SprintFunc()
	if statusCode >= 200 && statusCode < 300 {
		coloredStatus = color.New(color.BgGreen, color.FgBlack).SprintFunc()
	} else if statusCode >= 300 && statusCode < 400 {
		coloredStatus = color.New(color.BgYellow, color.FgBlack).SprintFunc()
	} else if statusCode >= 400 && statusCode < 600 {
		coloredStatus = color.New(color.BgRed, color.FgBlack).SprintFunc()
	}

	formatted := fmt.Sprintf(" %s ", status)
	return fmt.Sprintf("%s", coloredStatus(formatted))
}

func formatProtocol(protocol string) string {
	coloredProtocol := color.New(color.BgWhite, color.FgBlack).SprintFunc()
	formatted := fmt.Sprintf(" %s ", protocol)
	return fmt.Sprintf("%s", coloredProtocol(formatted))
}
