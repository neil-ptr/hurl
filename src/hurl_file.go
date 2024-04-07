package src

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"

	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"strings"
)

const (
	METHOD             = 0
	URL                = 1
	PROTOCOL           = 2
	NAME               = 0
	VALUE              = 1
	FILE_EMBED_TAG     = 0
	FILE_EMBED         = 1
	MULTIPART_FORM_TAG = 0
	MULTIPART_NAME     = 1
	MULTIPART_VALUE    = 2
	KEY                = 0
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

func extractFileEmbedPath(s string) string {
	filePathComponents := strings.SplitN(s, "=", 2)
	if len(filePathComponents) < 2 || len(filePathComponents) > 2 || filePathComponents[FILE_EMBED_TAG] != "@file" {
		return ""
	}

	return filePathComponents[FILE_EMBED]
}

func interpolateEnvVar(line []byte) (string, error) {
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

func parseKeyValPair(s string) (string, string, error) {
	kv := strings.Split(s, "=")
	if len(kv) != 2 {
		return "", "", fmt.Errorf("invalid key=value pair: %s\n", s)
	}

	key := kv[KEY]
	value := strings.Trim(kv[VALUE], "\"")

	return key, value, nil
}

func parseMultiPart(sc *bufio.Scanner) ([]MultiPartItem, error) {
	multipartItems := []MultiPartItem{}

	for sc.Scan() {
		multipartComponents := strings.Split(sc.Text(), ";")

		multipartTag := strings.TrimSpace(multipartComponents[MULTIPART_FORM_TAG])
		name := strings.TrimSpace(multipartComponents[MULTIPART_NAME])
		value := strings.TrimSpace(multipartComponents[MULTIPART_VALUE])

		if multipartTag != "form-data" {
			return []MultiPartItem{}, fmt.Errorf("found \"%s\" instead of \"form-data\" multipart tag", multipartTag)
		}

		formFieldNameKey, formFieldName, err := parseKeyValPair(name)
		if err != nil {
			return []MultiPartItem{}, err
		}
		if formFieldNameKey != "name" {
			return []MultiPartItem{}, fmt.Errorf("expected to find \"name\" in form data, found \"%s\"", formFieldNameKey)
		}

		formFieldValueKey, formFieldValue, err := parseKeyValPair(value)
		if err != nil {
			return []MultiPartItem{}, err
		}

		multipartItem := MultiPartItem{formFieldName, formFieldValueKey == "filename", formFieldValue}

		multipartItems = append(multipartItems, multipartItem)
	}

	return multipartItems, nil
}

func parseBody(sc *bufio.Scanner) ([]byte, bool, error) {
	bodyBuffer := bytes.Buffer{}

	fileEmbed := ""
	fileEmbedFound := false
	for sc.Scan() {
		line := sc.Text()

		extractedFileEmbed := extractFileEmbedPath(line)

		if fileEmbedFound && extractedFileEmbed != "" {
			return []byte{}, false, errors.New("more than 1 file embed found")
		}

		if len(extractedFileEmbed) > 0 {
			fileEmbedFound = true
			fileEmbed = extractedFileEmbed
		}

		bodyBuffer.Write(sc.Bytes())
		bodyBuffer.Write([]byte{'\n'}) // add newline because scanning removes it
	}

	if fileEmbedFound {
		if fileEmbed == "" {
			return []byte{}, false, errors.New("file embed path is empty")
		}

		return []byte(fileEmbed), true, nil
	}

	return bodyBuffer.Bytes(), false, nil
}

type MultiPartItem struct {
	Name       string
	IsFilePath bool
	Value      string
}

type HurlFile struct {
	Method            string
	URL               url.URL
	Headers           map[string]string
	Body              []byte
	FileEmbed         string
	MultipartFormData []MultiPartItem
	MultipartBoundary string

	// CLI and hurl.json options
	Config HurlConfig
}

func ParseHurlFile(r io.Reader) (*HurlFile, error) {

	h := &HurlFile{}
	sc := bufio.NewScanner(r)

	//=== request line ===//
	sc.Scan()
	requestLine := sc.Bytes()
	line, err := interpolateEnvVar(requestLine)
	if err != nil {
		return &HurlFile{}, err
	}

	requestLineComponents := strings.Split(string(line), " ")

	if len(requestLineComponents) > 2 {
		return &HurlFile{}, errors.New("Too many request line components")
	}
	if len(requestLineComponents) < 2 {
		return &HurlFile{}, errors.New("Not enough request line components")
	}

	parsedUrl, err := url.ParseRequestURI(requestLineComponents[URL])
	if err != nil {
		return &HurlFile{}, err
	}

	method := requestLineComponents[METHOD]
	if !isValidMethod(method) {
		return &HurlFile{}, errors.New("invalid HTTP method")
	}

	h.URL = *parsedUrl
	h.Method = requestLineComponents[METHOD]

	//=== headers ===//
	headerMap := make(map[string]string)

	scanFoundToken := sc.Scan()
	for scanFoundToken && strings.TrimSpace(sc.Text()) != "" {
		headerComponents := strings.SplitN(sc.Text(), ":", 2)
		if len(headerComponents) != 2 {
			return nil, fmt.Errorf("header is malformed: `%s`", sc.Text())
		}

		headerName, err := interpolateEnvVar([]byte(strings.TrimSpace(headerComponents[NAME])))
		if err != nil {
			return &HurlFile{}, fmt.Errorf("error interpolating value")
		}

		headerVal, err := interpolateEnvVar([]byte(strings.TrimSpace(headerComponents[VALUE])))
		if err != nil {
			return &HurlFile{}, fmt.Errorf("error interpolating value")
		}

		headerMap[headerName] = headerVal

		scanFoundToken = sc.Scan()
	}
	if sc.Err() != nil {
		return &HurlFile{}, sc.Err()
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

	//=== body ===//

	// TODO: check valid content-type
	headerContentType, exists := h.Headers["Content-Type"]
	if exists && headerContentType == "multipart/form-data" {
		// form data
		multipartFormData, err := parseMultiPart(sc)
		if err != nil {
			return &HurlFile{}, err
		}

		h.MultipartFormData = multipartFormData

	} else if exists && headerContentType != "multipart/form-data" {
		// read body as is, might have file embed
		body, containsFileEmbed, err := parseBody(sc)
		if err != nil {
			return &HurlFile{}, err
		}

		if containsFileEmbed {
			h.FileEmbed = string(body)
		} else {
			h.Body = body
		}
	} else {
		err := errors.New("no \"Content-Type\" header found, using \"text/plain\" as \"Content-Type\" header")
		h.Headers["Content-Type"] = "text/plain"
		PrintWarning(err)

		body, _, err := parseBody(sc)
		if err != nil {
			return &HurlFile{}, err
		}
		h.Body = body
	}

	return h, nil
}

func WriteMultipart(h HurlFile, writer *multipart.Writer) error {
	for _, multiPartItem := range h.MultipartFormData {
		if multiPartItem.IsFilePath {
			file, err := os.Stat(multiPartItem.Value)
			if err != nil {
				return nil
			}

			fileData, err := os.ReadFile(multiPartItem.Value)
			if err != nil {
				return nil
			}

			mimeHeader := make(textproto.MIMEHeader)

			mimeHeader.Set("Content-Disposition", fmt.Sprintf("form-data; name=\"%s\"; filename=\"%s\"", multiPartItem.Name, file.Name()))
			mimeHeader.Set("Content-Type", http.DetectContentType(fileData))

			part, err := writer.CreatePart(mimeHeader)
			if err != nil {
				return err
			}

			part.Write(fileData)
		} else {
			err := writer.WriteField(multiPartItem.Name, multiPartItem.Value)
			if err != nil {
				return err
			}
		}
	}

	writer.Close()
	return nil
}

func (h *HurlFile) NewRequest() (*http.Request, error) {
	body := &bytes.Buffer{}

	contentType, exists := h.Headers["Content-Type"]
	if !exists && h.Method != "GET" && h.Method != "DELETE" {
		return &http.Request{}, errors.New("\"Content-Type\" header missing")
	}

	if contentType == "multipart/form-data" && len(h.MultipartFormData) > 0 {
		writer := multipart.NewWriter(body)
		h.Headers["Content-Type"] = fmt.Sprintf("%s; boundary=%s", contentType, writer.Boundary())
		h.MultipartBoundary = writer.Boundary()

		WriteMultipart(*h, writer)
	} else if h.FileEmbed != "" {
		fileData, err := os.ReadFile(h.FileEmbed)
		if err != nil {
			return &http.Request{}, nil
		}
		body.Write(fileData)
	} else {
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
