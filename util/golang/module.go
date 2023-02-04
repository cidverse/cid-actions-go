package golang

import (
	"strings"

	cidsdk "github.com/cidverse/cid-sdk-go"
)

func IsGoLibrary(module *cidsdk.ProjectModule) bool {
	for _, path := range module.Files {
		file := strings.TrimPrefix(strings.TrimPrefix(path, module.ModuleDir), "/")

		if !strings.Contains(file, "/") && strings.HasSuffix(file, ".go") {
			return false
		}
	}

	return true
}
