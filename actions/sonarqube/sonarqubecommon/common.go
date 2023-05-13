package sonarqubecommon

import (
	cidsdk "github.com/cidverse/cid-sdk-go"
)

func SonarqubeGoModTestData(debug bool) cidsdk.ProjectActionData {
	return cidsdk.ProjectActionData{
		ProjectDir: "/my-project",
		Config: cidsdk.CurrentConfig{
			Debug:       debug,
			Log:         map[string]string{},
			ProjectDir:  "/my-project",
			ArtifactDir: "/my-project/.dist",
			TempDir:     "/my-project/.tmp",
		},
		Modules: []cidsdk.ProjectModule{
			{
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
		},
		Env: map[string]string{
			"NCI_PROJECT_NAME":        "my-project-name",
			"NCI_PROJECT_DESCRIPTION": "my-project-description",
		},
	}
}
