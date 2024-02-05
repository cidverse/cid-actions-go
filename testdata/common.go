package testdata

import (
	cidsdk "github.com/cidverse/cid-sdk-go"
)

func ModuleDefault(env map[string]string, debug bool) cidsdk.ProjectActionData {
	m := cidsdk.ProjectActionData{
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
			"NCI_REPOSITORY_KIND":        "git",
			"NCI_REPOSITORY_REMOTE":      "https://github.com/cidverse/normalizeci.git",
			"NCI_REPOSITORY_URL":         "https://github.com/cidverse/normalizeci",
			"NCI_REPOSITORY_HOST_SERVER": "github.com",
			"NCI_COMMIT_REF_NAME":        "v1.2.0",
			"NCI_COMMIT_HASH":            "abcdef123456",
			"NCI_COMMIT_REF_VCS":         "refs/tags/v1.2.0",
			"NCI_PROJECT_URL":            "https://github.com/cidverse/normalizeci",
		},
	}
	for k, v := range env {
		m.Env[k] = v
	}
	return m
}
