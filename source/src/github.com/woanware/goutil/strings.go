package goutil

import (
	"strings"
	"fmt"
)

// Emulates the python partition function
func Partition(data string, separator string) (pre string, post string) {
	index := strings.Index(data, separator)
	if index == -1 {
		return "", ""
	}

	return data[:index], data[index+1:]
}

//
func GetSeparator(s string) rune {
	var sep string
	s = `"` + s + `"`
	fmt.Sscanf(s, "%q", &sep)

	return ([]rune(sep))[0]
}

// Removes the leading/trailing quotes
func RemoveQuotes(data string) string {
	data = strings.TrimSpace(data)
	if len(data) == 0 {
		return ""
	}

	if data == "\"\"" {
		return ""
	}

	if data[:1] == "\"" {
		data = data[1:len(data) - 1]
	}

	if data[len(data)-1:] == "\"" {
		data = data[:len(data)-1]
	}

	return data
}
