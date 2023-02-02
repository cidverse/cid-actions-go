package container

import (
	"strings"
)

func GetImageReferenceWithoutTag(input string) string {
	parts := strings.SplitN(input, ":", 2)
	return parts[0]
}
