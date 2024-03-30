package src

import "fmt"

func InterpolateValues(s string, values map[string]string) string {
	for _, c := range s {
		fmt.Println(byte(c))
	}

	return ""
}
