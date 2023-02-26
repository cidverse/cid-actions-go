package cosign

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

func DockerfileTestData(debug bool) cidsdk.ModuleActionData {
	return cidsdk.ModuleActionData{
		ProjectDir: "/my-project",
		Module: cidsdk.ProjectModule{
			ProjectDir:        "/my-project",
			ModuleDir:         "/my-project",
			Discovery:         []cidsdk.ProjectModuleDiscovery{{File: "/my-project/Dockerfile"}},
			Name:              "my-project",
			Slug:              "my-project",
			BuildSystem:       string(cidsdk.BuildSystemContainer),
			BuildSystemSyntax: string(cidsdk.ContainerFile),
			Language:          &map[string]string{},
			Submodules:        nil,
		},
		Config: cidsdk.CurrentConfig{
			Debug:       debug,
			Log:         map[string]string{},
			ProjectDir:  "/my-project",
			ArtifactDir: "/my-project/.dist",
			TempDir:     "/my-project/.tmp",
		},
		Env: map[string]string{
			"NCI_CONTAINERREGISTRY_HOST":       "ghcr.io",
			"NCI_CONTAINERREGISTRY_REPOSITORY": "cidverse/dummy",
			"NCI_CONTAINERREGISTRY_TAG":        "latest",
			"NCI_REPOSITORY_KIND":              "git",
			"NCI_REPOSITORY_REMOTE":            "https://github.com/cidverse/normalizeci.git",
			"NCI_COMMIT_SHA":                   "abcdef123456",
		},
	}
}
