package python

import (
	"os"
	"testing"

	cidsdk "github.com/cidverse/cid-sdk-go"
)

func TestMain(m *testing.M) {
	cidsdk.JoinSeparator = "/"
	code := m.Run()
	os.Exit(code)
}

func PythonTestData(buildSystem string, debug bool) cidsdk.ModuleActionData {
	return cidsdk.ModuleActionData{
		Module: cidsdk.ProjectModule{
			ProjectDir:        "/my-project",
			ModuleDir:         "/my-project",
			Discovery:         []cidsdk.ProjectModuleDiscovery{{File: "/my-project/package.json"}},
			Name:              "my-package",
			Slug:              "my-package",
			BuildSystem:       buildSystem,
			BuildSystemSyntax: "default",
			Language:          &map[string]string{},
			Submodules:        nil,
		},
		Config: cidsdk.CurrentConfig{
			Debug:       debug,
			Log:         map[string]string{},
			ArtifactDir: ".dist",
			TempDir:     ".tmp",
		},
	}
}
