package testdata

import (
	cidsdk "github.com/cidverse/cid-sdk-go"
)

func ModuleDockerfile() cidsdk.ModuleActionData {
	return cidsdk.ModuleActionData{
		ProjectDir: "/my-project",
		Module: cidsdk.ProjectModule{
			ProjectDir:        "/my-project",
			ModuleDir:         "/my-project",
			Discovery:         []cidsdk.ProjectModuleDiscovery{{File: "/my-project/Dockerfile"}},
			Name:              "my-project",
			Slug:              "my-project",
			BuildSystem:       string(cidsdk.BuildSystemContainer),
			BuildSystemSyntax: string(cidsdk.BuildSystemSyntaxContainerFile),
			Language:          &map[string]string{},
			Submodules:        nil,
		},
		Config: cidsdk.CurrentConfig{
			Debug:       false,
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
			"NCI_COMMIT_HASH":                  "abcdef123456",
		},
	}
}
