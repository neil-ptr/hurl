package src

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessLineSuccessNoTemplate(t *testing.T) {
	l := []byte("what the flip")
	line, err := processLine(l)
	fmt.Println(line)

	assert.Nil(t, err)

	assert.Equal(t, string(l), line)
}

func TestProcessLineSuccessTemplateEndingEdgeCase(t *testing.T) {
	baseUrl := "https://jsonplaceholder.typicode.com"
	t.Setenv("BASE_URL", baseUrl)

	line := []byte("{{BASE_URL}}")
	answer := fmt.Sprintf("%s", baseUrl)

	processedLine, err := processLine(line)

	assert.Nil(t, err)

	assert.Equal(t, answer, processedLine)
}

func TestProcessLineSuccessTemplate(t *testing.T) {
	baseUrl := "https://jsonplaceholder.typicode.com"
	t.Setenv("BASE_URL", baseUrl)

	line := []byte("GET {{BASE_URL}}/todos/1")
	answer := fmt.Sprintf("GET %s/todos/1", baseUrl)

	processedLine, err := processLine(line)

	assert.Nil(t, err)

	assert.Equal(t, answer, processedLine)
}

func TestProcessLineSuccessTemplateSpacesAndTabs(t *testing.T) {
	baseUrl := "https://jsonplaceholder.typicode.com"
	t.Setenv("BASE_URL", baseUrl)

	line := []byte("GET {{				BASE_URL   }}/todos/1")
	answer := fmt.Sprintf("GET %s/todos/1", baseUrl)

	processedLine, err := processLine(line)

	assert.Nil(t, err)

	assert.Equal(t, answer, processedLine)
}

func TestProcessLineFailureInvalidCharacter(t *testing.T) {
	baseUrl := "https://jsonplaceholder.typicode.com"
	t.Setenv("BASE_URL", baseUrl)

	line := []byte("GET {{B%%SE_URL}}/todos/2")

	_, err := processLine(line)

	assert.ErrorContains(t, err, "template variable contains invalid character")
}
func TestProcessLineFailureInvalidFirstChar(t *testing.T) {
	line := []byte("GET {{1BASE_URL}}/todos/2")

	_, err := processLine(line)

	assert.ErrorContains(t, err, "template variable must begin with letter")
}

func TestProcessLineFailureEmptyTemplateVar(t *testing.T) {
	line := []byte("GET {{}}/todos/2")

	_, err := processLine(line)

	assert.ErrorContains(t, err, "template variable cannot be empty")
}

func TestParseHurlFileNoBody(t *testing.T) {
	r := strings.NewReader("GET https://example.com")

	parsedUrl, _ := url.Parse("https://example.com")

	hurlFile, err := ParseHurlFile(r)

	assert.Nil(t, err)

	assert.Equal(t, *parsedUrl, hurlFile.URL)
}

func TestParseHurlFileReadableBody(t *testing.T) {
	r := strings.NewReader("POST https://example.com\nContent-Type: application/json\n\n{\"hi\": 1}")

	parsedUrl, _ := url.Parse("https://example.com")

	headers := make(map[string]string)
	headers["User-Agent"] = "hurl/0.1.0"
	headers["Content-Type"] = "application/json"

	hurlFile, err := ParseHurlFile(r)

	assert.Nil(t, err)

	assert.Equal(t, *parsedUrl, hurlFile.URL)
	assert.Equal(t, headers, hurlFile.Headers)
}

func TestParseHurlFileFilePaths(t *testing.T) {
	r := strings.NewReader("POST https://example.com\nContent-Type: image/png\n\n@file=path/idk.png")

	parsedUrl, _ := url.Parse("https://example.com")

	headers := make(map[string]string)
	headers["User-Agent"] = "hurl/0.1.0"
	headers["Content-Type"] = "image/png"

	hurlFile, err := ParseHurlFile(r)

	assert.Nil(t, err)

	assert.Equal(t, *parsedUrl, hurlFile.URL)
	assert.Equal(t, headers, hurlFile.Headers)
	assert.Equal(t, []string{"path/idk.png"}, hurlFile.FilePaths)
}
