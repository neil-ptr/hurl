package src

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const (
	METHOD   = 0
	URL      = 1
	PROTOCOL = 2
	NAME     = 0
	VALUE    = 1
)

func isNum(c byte) bool {
	return 48 <= c && c <= 57
}

func isUnderscore(c byte) bool {
	return c == 95 // 95 is the "_" character in ascii
}

func isAlpha(c byte) bool {
	return (65 <= c && c <= 90) || (97 <= c && c <= 122)
}

func validateTemplateVariable(v []byte) error {
	if len(v) == 0 {
		return errors.New("template variable cannot be empty")
	}

	if !isAlpha(v[0]) {
		return fmt.Errorf("template variable must begin with letter: %s", string(v))
	}

	for _, c := range v {
		if isNum(c) || isAlpha(c) || isUnderscore(c) {
			continue
		}

		return fmt.Errorf("template variable contains invalid character: %s", string(v))
	}

	return nil
}

func isValidMethod(m string) bool {
	return m == "GET" || m == "POST" || m == "PUT" || m == "PATCH" || m == "DELETE"
}

func isFilePath(s string) bool {
	filePathComponents := strings.SplitN(s, "=", 2)
	if len(filePathComponents) < 2 || len(filePathComponents) < 2 || filePathComponents[0] != "@file" {
		return false
	}

	return true
}

func processLine(line []byte) (string, error) {
	processedLine := []byte{}
	trimmed := bytes.TrimSpace(line)

	environmentVariableSet := make(map[string]void)

	i := 0
	for i < len(trimmed) {
		templateVar := []byte{}

		for i < len(trimmed)-1 && string(trimmed[i:i+2]) == "{{" {

			// skip to just after the second opening brace
			i += 2

			for i < len(trimmed)-1 && string(trimmed[i:i+2]) != "}}" {
				templateVar = append(templateVar, trimmed[i])
				i++
			}

			trimmedTemplateVar := bytes.TrimSpace(templateVar)
			err := validateTemplateVariable(trimmedTemplateVar)
			if err != nil {
				return "", err
			}

			envVar := os.Getenv(string(trimmedTemplateVar))
			if len(envVar) > 0 {
				processedLine = append(processedLine, envVar...)
				environmentVariableSet[string(trimmedTemplateVar)] = member
			} else {
				warning := fmt.Errorf("warning: could not find environment variable: %s\n", trimmedTemplateVar)
				PrintWarning(warning)
			}

			// skip to just after second closing brace
			i += 2

		}

		if i < len(trimmed) {
			processedLine = append(processedLine, trimmed[i])
		}

		i++
	}

	return string(processedLine), nil
}

func isHumanReadableContentType(contentType string) bool {
	return contentType == "application/json" || contentType == "text/raw" || contentType == "text/html"
}

func parseHumanReadableBody(sc *bufio.Scanner) ([]byte, error) {
	body := []byte{}
	for sc.Scan() {
		body = append(body, sc.Bytes()...)

		// scanning removes newline I guess, add it back to make sure it looks
		// the same in the hurl file as it looks in the output for consistency
		body = append(body, '\n')
	}

	if sc.Err() != nil {
		return []byte{}, sc.Err()
	}

	return body, nil
}

func parseFilePaths(body string) ([]string, error) {
	filePaths := []string{}

	lines := strings.Split(body, "\n")
	for _, line := range lines {
		// filePathLine can only be 2 elements long and
		// must have an empty first element after splitting on file path keyword
		if isFilePath(line) {
			return []string{}, errors.New("invalid file path")
		}

		filePathLine := strings.SplitN(line, "=", 2)
		filePath := filePathLine[1]

		filePaths = append(filePaths, filePath)
	}

	return filePaths, nil
}

type HurlFile struct {
	Method    string
	URL       url.URL
	Headers   map[string]string
	Body      []byte
	FilePaths []string

	// CLI and hurlrc options
	Config HurlConfig
}

func ParseHurlFile(r io.Reader) (HurlFile, error) {
	h := HurlFile{}
	sc := bufio.NewScanner(r)

	// request line
	sc.Scan()
	requestLine := sc.Bytes()
	line, err := processLine(requestLine)
	if err != nil {
		return HurlFile{}, err
	}

	requestLineComponents := strings.Split(string(line), " ")

	if len(requestLineComponents) > 2 {
		return HurlFile{}, errors.New("Too many request line components")
	}
	if len(requestLineComponents) < 2 {
		return HurlFile{}, errors.New("Not enough request line components")
	}

	parsedUrl, err := url.ParseRequestURI(requestLineComponents[URL])
	if err != nil {
		return HurlFile{}, err
	}

	method := requestLineComponents[METHOD]
	if !isValidMethod(method) {
		return HurlFile{}, errors.New("invalid HTTP method")
	}

	h.URL = *parsedUrl
	h.Method = requestLineComponents[METHOD]

	// headers
	headerMap := make(map[string]string)

	scanFoundToken := sc.Scan()
	for scanFoundToken && strings.TrimSpace(sc.Text()) != "" {
		headerComponents := strings.Split(sc.Text(), ": ")
		headerName := headerComponents[NAME]
		headerVal := headerComponents[VALUE]

		headerMap[headerName] = headerVal

		scanFoundToken = sc.Scan()
	}
	if sc.Err() != nil {
		return HurlFile{}, sc.Err()
	}

	h.Headers = headerMap

	hostHeaderVal, exists := h.Headers["Host"]
	if exists && h.URL.Hostname() != hostHeaderVal {
		PrintWarning(errors.New("host header value does not match host in URL, using host in URL"))
		h.Headers["Host"] = h.URL.Hostname()
	}

	h.Headers["User-Agent"] = "hurl/0.1.0"

	// no body, done
	if !scanFoundToken {
		return h, nil
	}

	// body
	filePathPresent := false
	bodyBuffer := bytes.Buffer{}
	for sc.Scan() {
		line := sc.Text()
		if isFilePath(line) {
			filePathPresent = true
		}

		bodyBuffer.Write(sc.Bytes())

		// add newline because scanning removes it
		bodyBuffer.Write([]byte{'\n'})
	}
	if sc.Err() != nil {
		return HurlFile{}, sc.Err()
	}

	if filePathPresent {
		// read file paths
		filePaths, err := parseFilePaths(bodyBuffer.String())
		if err != nil {
			return HurlFile{}, err
		}

		h.FilePaths = filePaths
		fmt.Println(filePaths)

	} else {
		_, exists := h.Headers["Content-Type"]
		if !exists {
			err := errors.New("no \"Content-Type\" defined, using \"text/plain\" content type")
			PrintWarning(err)
		}

		// read as raw text
		h.Body = bodyBuffer.Bytes()
	}

	return h, nil
}

func (h HurlFile) NewRequest() (*http.Request, error) {
	body := &bytes.Buffer{}

	for _, filePath := range h.FilePaths {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}

		defer file.Close()

		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", filepath.Base(file.Name()))
		io.Copy(part, file)
		writer.Close()
	}

	if len(h.Body) > 0 {
		body.Write(h.Body)
	}

	req, err := http.NewRequest(h.Method, h.URL.String(), body)
	if err != nil {
		return nil, err
	}

	header := http.Header{}
	for name, val := range h.Headers {
		header[name] = []string{val}
	}
	req.Header = header

	return req, nil
}
