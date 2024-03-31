package src

import (
	"github.com/fatih/color"
)

// Set data structure helpers
type void struct{}

var member void

func PrintWarning(err error) {
	warning := color.New(color.Bold, color.FgRed).PrintfFunc()
	warning("warning: %s\n", err.Error())
}
