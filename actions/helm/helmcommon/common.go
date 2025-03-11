package helmcommon

import (
	cidsdk "github.com/cidverse/cid-sdk-go"
)

const HelmVersionConstraint = ">=3.16.0"

func GetHelmTestData(env map[string]string, debug bool) *cidsdk.ModuleActionData {
	return &cidsdk.ModuleActionData{
		ProjectDir: "/my-project",
		Module: cidsdk.ProjectModule{
			ProjectDir:        "/my-project",
			ModuleDir:         "/my-project/charts/mychart",
			Discovery:         []cidsdk.ProjectModuleDiscovery{{File: "/my-project/package.json"}},
			Name:              "my-package",
			Slug:              "my-package",
			BuildSystem:       string(cidsdk.BuildSystemHelm),
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
		Env: env,
	}
}
