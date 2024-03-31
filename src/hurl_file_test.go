package src

import (
	"fmt"
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

	assert.ErrorContains(t, err, "template variable has to begin with letter")
}

func TestProcessLineFailureEmptyTemplateVar(t *testing.T) {
	line := []byte("GET {{}}/todos/2")

	_, err := processLine(line)

	assert.ErrorContains(t, err, "template variable cannot be empty")
}

func TestParseHurlFile(t *testing.T) {

}
