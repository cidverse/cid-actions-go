package testdata

import (
	cidsdk "github.com/cidverse/cid-sdk-go"
)

func ModuleRenovate() cidsdk.ModuleActionData {
	return cidsdk.ModuleActionData{
		Module: cidsdk.ProjectModule{
			ProjectDir:        "/my-project",
			ModuleDir:         "/my-project",
			Discovery:         []cidsdk.ProjectModuleDiscovery{{File: "/my-project/renovate.json"}},
			Name:              "my-package",
			Slug:              "my-package",
			Type:              "specification",
			SpecificationType: "renovate",
			Language:          &map[string]string{},
			Submodules:        nil,
		},
		Config: cidsdk.CurrentConfig{
			Debug:       false,
			Log:         map[string]string{},
			ArtifactDir: ".dist",
			TempDir:     ".tmp",
		},
	}
}
