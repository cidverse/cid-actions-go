package java

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

func GradleTestData(debug bool) cidsdk.ModuleActionData {
	return cidsdk.ModuleActionData{
		ProjectDir: "/my-project",
		Module: cidsdk.ProjectModule{
			ProjectDir:        "/my-project",
			ModuleDir:         "/my-project",
			Discovery:         []cidsdk.ProjectModuleDiscovery{{File: "/my-project/build.gradle.kts"}},
			Name:              "my-module",
			Slug:              "my-module",
			BuildSystem:       string(cidsdk.BuildSystemGradle),
			BuildSystemSyntax: string(cidsdk.GradleKotlinDSL),
			Language:          &map[string]string{},
			Submodules:        nil,
		},
		Config: cidsdk.CurrentConfig{
			Debug:       debug,
			Log:         map[string]string{},
			ArtifactDir: "/my-project/.dist",
			TempDir:     "/my-project/.tmp",
		},
		Env: map[string]string{
			"NCI_COMMIT_REF_TYPE":    "branch",
			"NCI_COMMIT_REF_RELEASE": "main",
		},
	}
}
