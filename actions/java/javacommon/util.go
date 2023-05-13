package javacommon

import (
	"strings"
)

// GetVersion returns the suggested java artifact version
// Unless the reference is a git tag versions will get a -SNAPSHOT suffix
func GetVersion(refType string, refName string, shortHash string) string {
	if refType == "tag" {
		return refName
	}

	return strings.TrimLeft(refName, "v") + "-" + shortHash
}
