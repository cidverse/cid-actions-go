package helm

import cidsdk "github.com/cidverse/cid-sdk-go"

func GetHelmTestData(debug bool) cidsdk.ModuleActionData {
	return cidsdk.ModuleActionData{
		ProjectDir: "/my-project",
		Module: cidsdk.ProjectModule{
			ProjectDir:        "/my-project",
			ModuleDir:         "/my-project",
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
	}
}
