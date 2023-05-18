package api

import (
	cidsdk "github.com/cidverse/cid-sdk-go"
)

func GetProjectActionData(debug bool) cidsdk.ProjectActionData {
	return cidsdk.ProjectActionData{
		ProjectDir: "/my-project",
		Config: cidsdk.CurrentConfig{
			Debug:       debug,
			Log:         map[string]string{},
			ProjectDir:  "/my-project",
			ArtifactDir: "/my-project/.dist",
			TempDir:     "/my-project/.tmp",
		},
		Modules: nil,
		Env: map[string]string{
			"NCI_REPOSITORY_KIND":   "git",
			"NCI_REPOSITORY_REMOTE": "https://github.com/cidverse/normalizeci.git",
			"NCI_REPOSITORY_URL":    "https://github.com/cidverse/normalizeci",
			"NCI_COMMIT_REF_NAME":   "v1.2.0",
			"NCI_COMMIT_HASH":       "abcdef123456",
			"NCI_COMMIT_REF_VCS":    "refs/tags/v1.2.0",
		},
	}
}

func GetUnknownTestData(debug bool) cidsdk.ModuleActionData {
	return cidsdk.ModuleActionData{
		Module: cidsdk.ProjectModule{
			ProjectDir:        "/my-project",
			ModuleDir:         "/my-project",
			Discovery:         []cidsdk.ProjectModuleDiscovery{{File: "/my-project/go.mod"}},
			Name:              "github.com/cidverse/my-project",
			Slug:              "github-com-cidverse-my-project",
			BuildSystem:       "unknown",
			BuildSystemSyntax: "default",
			Language:          &map[string]string{"go": "1.19.0"},
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

func GetNodeTestData(debug bool) cidsdk.ModuleActionData {
	return cidsdk.ModuleActionData{
		Module: cidsdk.ProjectModule{
			ProjectDir:        "/my-project",
			ModuleDir:         "/my-project",
			Discovery:         []cidsdk.ProjectModuleDiscovery{{File: "/my-project/package.json"}},
			Name:              "my-package",
			Slug:              "my-package",
			BuildSystem:       "node",
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
			"NCI_COMMIT_REF_TYPE":     "branch",
			"NCI_COMMIT_REF_NAME":     "main",
		},
	}
}
