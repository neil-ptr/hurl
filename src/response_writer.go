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

func highlight(body string) string {
	return ""
}

func WriteResponse(res *http.Response) {
	coloredStatus := color.New(color.FgWhite, color.Bold).SprintFunc()
	if res.StatusCode >= 200 && res.StatusCode < 300 {
		coloredStatus = color.New(color.FgGreen, color.Bold).SprintFunc()
	} else if res.StatusCode >= 300 && res.StatusCode < 400 {
		coloredStatus = color.New(color.FgYellow, color.Bold).SprintFunc()
	} else if res.StatusCode >= 400 && res.StatusCode < 600 {
		coloredStatus = color.New(color.FgRed, color.Bold).SprintFunc()
	}

	fmt.Printf("< %s %s\n", res.Proto, coloredStatus(res.Status))

	coloredHeaderName := color.New(color.FgYellow).SprintFunc()
	for name, value := range res.Header {
		fmt.Printf("< %s: %s\n", coloredHeaderName(name), strings.Join(value, ""))

	}

	// separate headers from body
	fmt.Printf("\n")

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("hurl: ", err)
		return
	}

	buffer := bytes.Buffer{}

	err = quick.Highlight(&buffer, string(body), "json", "terminal256", "")

	fmt.Println(buffer.String())
}
