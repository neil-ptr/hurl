package src

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
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

func formatPath(p string) string {
	coloredPath := color.New(color.BgWhite, color.FgBlack).SprintFunc()
	correctedPath := p
	if len(p) == 0 {
		correctedPath = "/"
	}
	formatted := fmt.Sprintf(" %s ", correctedPath)
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

func PrettifyJson(j []byte) ([]byte, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, j, "", "  "); err != nil {
		return []byte{}, err
	}

	return prettyJSON.Bytes(), nil
}

func FormatHeaders(headers http.Header, directionCharacter string) []byte {
	buffer := bytes.Buffer{}

	for name, value := range headers {
		yellow := color.New(color.FgYellow).SprintFunc()
		headerVal := strings.Join(value, "")

		formattedHeader := fmt.Sprintf("%s %s: %s\n", directionCharacter, yellow(name), headerVal)

		buffer.Write([]byte(formattedHeader))
	}

	return buffer.Bytes()
}

func FormatFilePathsTitle() string {
	title := color.New(color.FgBlack, color.BgWhite).SprintFunc()
	return fmt.Sprintf("%s\n", title(" body contents outputted to: "))
}

func FormatFileEmbed(fileEmbed string) []byte {
	buffer := bytes.Buffer{}

	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	formattedFilePath := fmt.Sprintf("%s%s=%s", red("@"), green("file"), fileEmbed)

	buffer.Write([]byte(formattedFilePath))

	return buffer.Bytes()
}

func FormatBody(body []byte, mediaType string) ([]byte, error) {
	buffer := bytes.Buffer{}

	formatter := "" // raw text
	if mediaType == "application/json" {
		formatter = "json"
		prettified, err := PrettifyJson(body)
		if err != nil {
			return []byte{}, nil
		}

		body = prettified
	} else if mediaType == "text/html" {
		formatter = "html"
	}

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

func FormatMultiPart(multipartItems []MultiPartItem, boundary string) ([]byte, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.SetBoundary(boundary)

	for _, multiPartItem := range multipartItems {
		if multiPartItem.IsFilePath {
			file, err := os.Stat(multiPartItem.Value)
			if err != nil {
				return []byte{}, err
			}

			fileData, err := os.ReadFile(multiPartItem.Value)
			if err != nil {
				return []byte{}, err
			}

			mimeHeader := make(textproto.MIMEHeader)

			mimeHeader.Set("Content-Disposition", fmt.Sprintf("form-data; name=\"%s\"; filename=\"%s\"", multiPartItem.Name, file.Name()))
			mimeHeader.Set("Content-Type", http.DetectContentType(fileData))

			part, err := writer.CreatePart(mimeHeader)
			if err != nil {
				return []byte{}, err
			}

			part.Write(FormatFileEmbed(multiPartItem.Value))
		} else {
			err := writer.WriteField(multiPartItem.Name, multiPartItem.Value)
			if err != nil {
				return []byte{}, err
			}
		}
	}

	writer.Close()

	return body.Bytes(), nil
}
