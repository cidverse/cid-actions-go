package qodanascancommon

import (
	cidsdk "github.com/cidverse/cid-sdk-go"
)

func GoModuleTestData() *cidsdk.ModuleActionData {
	return &cidsdk.ModuleActionData{
		ProjectDir: "/my-project",
		Module: cidsdk.ProjectModule{
			ProjectDir:        "/my-project",
			ModuleDir:         "/my-project",
			Discovery:         []cidsdk.ProjectModuleDiscovery{{File: "/my-project/go.mod"}},
			Name:              "github.com/cidverse/my-project",
			Slug:              "github-com-cidverse-my-project",
			BuildSystem:       "gomod",
			BuildSystemSyntax: "default",
			Language:          &map[string]string{"go": "1.19.0"},
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
			"QODANA_EAP": "true",
		},
	}
}
