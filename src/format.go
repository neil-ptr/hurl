package src

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/alecthomas/chroma/v2/quick"
	"github.com/fatih/color"
)

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

func formatProtocol(protocol string) string {
	coloredProtocol := color.New(color.BgWhite, color.FgBlack).SprintFunc()
	formatted := fmt.Sprintf(" %s ", protocol)
	return fmt.Sprintf("%s", coloredProtocol(formatted))
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

func FormatHeaders(headers http.Header) []byte {
	buffer := bytes.Buffer{}

	for name, value := range headers {
		yellow := color.New(color.FgYellow).SprintFunc()
		headerVal := strings.Join(value, "")

		formattedHeader := fmt.Sprintf("> %s: %s\n", yellow(name), headerVal)

		buffer.Write([]byte(formattedHeader))
	}

	return buffer.Bytes()
}

func FormatFilePathsTitle() string {
	title := color.New(color.FgBlack, color.BgWhite).SprintFunc()
	return fmt.Sprintf("%s\n", title(" body contents outputted to: "))
}

func FormatFilePaths(filePaths []string) []byte {
	buffer := bytes.Buffer{}

	for _, filePath := range filePaths {

		red := color.New(color.FgRed).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()

		formattedFilePath := fmt.Sprintf("%s%s=%s\n", red("@"), green("file"), filePath)

		buffer.Write([]byte(formattedFilePath))
	}

	return buffer.Bytes()
}

func FormatBody(body []byte, contentType string) ([]byte, error) {
	formatter := "" // raw text
	if strings.Contains(contentType, "application/json") {
		formatter = "json"
	} else if strings.Contains(contentType, "text/html") {
		formatter = "html"
	}

	if len(body) == 0 {
		return []byte{}, nil
	}

	buffer := bytes.Buffer{}
	err := quick.Highlight(&buffer, string(body), formatter, "terminal", "")
	if err != nil {
		return []byte{}, nil
	}

	return buffer.Bytes(), nil
}

func FormatRequestLine(req http.Request) []byte {
	method := formatMethod(req.Method)
	path := formatPath(req.URL.Path)
	protocol := formatProtocol(req.Proto)

	formattedStatusline := fmt.Sprintf("> %s%s%s\n", method, path, protocol)

	return []byte(formattedStatusline)
}

func FormatStatusLine(res http.Response) string {
	protocol := formatProtocol(res.Proto)
	status := formatStatusCode(res.StatusCode, res.Status)

	return fmt.Sprintf("< %s%s\n", protocol, status)
}
