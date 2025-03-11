package util

import (
	"strings"
	"unicode"
)

func FirstNonEmpty(strings []string) string {
	for _, str := range strings {
		if str != "" {
			return str
		}
	}
	return ""
}

func TrimLeftEachLine(s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimLeftFunc(line, unicode.IsSpace)
	}
	return strings.Join(lines, "\n")
}
