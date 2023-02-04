package golang

import (
	"testing"

	cidsdk "github.com/cidverse/cid-sdk-go"
)

func TestIsGoLibrary(t *testing.T) {
	tests := []struct {
		module    cidsdk.ProjectModule
		isLibrary bool
	}{
		{
			module: cidsdk.ProjectModule{
				ModuleDir: "/path/to/module",
				Files: []string{
					"/path/to/module/main.go",
					"/path/to/module/subdir/subfile.go",
				},
			},
			isLibrary: false,
		},
		{
			module: cidsdk.ProjectModule{
				ModuleDir: "/path/to/module",
				Files: []string{
					"/path/to/module/subdir/subfile.go",
				},
			},
			isLibrary: true,
		},
	}

	for i, test := range tests {
		result := IsGoLibrary(&test.module)
		if result != test.isLibrary {
			t.Errorf("Test case %d failed: expected %t, got %t", i+1, test.isLibrary, result)
		}
	}
}
