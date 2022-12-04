package techdocs

import cidsdk "github.com/cidverse/cid-sdk-go"

func TechdocsTestData(debug bool) cidsdk.ModuleActionData {
	return cidsdk.ModuleActionData{
		ProjectDir: "/my-project",
		Module: cidsdk.ProjectModule{
			ProjectDir:        "/my-project",
			ModuleDir:         "/my-project/docs",
			Discovery:         []cidsdk.ProjectModuleDiscovery{{File: "/my-project/mkdocs.yaml"}},
			Name:              "my-module",
			Slug:              "my-module",
			BuildSystem:       string(cidsdk.BuildSystemMkdocs),
			BuildSystemSyntax: string(cidsdk.MkdocsTechdocs),
			Language:          &map[string]string{},
			Submodules:        nil,
		},
		Config: cidsdk.CurrentConfig{
			Debug:       debug,
			Log:         map[string]string{},
			ArtifactDir: "/my-project/.dist",
			TempDir:     "/my-project/.tmp",
		},
	}
}
