package testdata

import (
	cidsdk "github.com/cidverse/cid-sdk-go"
)

func ProjectRenovate() *cidsdk.ProjectActionData {
	return &cidsdk.ProjectActionData{
		ProjectDir: "/my-project",
		Config: cidsdk.CurrentConfig{
			Debug:       false,
			Log:         map[string]string{},
			ArtifactDir: ".dist",
			TempDir:     ".tmp",
		},
	}
}
