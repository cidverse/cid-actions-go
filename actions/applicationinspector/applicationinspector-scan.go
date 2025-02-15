package applicationinspector

import (
	"fmt"
	"strings"

	cidsdk "github.com/cidverse/cid-sdk-go"
)

type Action struct {
	Sdk cidsdk.SDKClient
}

type Config struct {
}

func (a Action) Metadata() cidsdk.ActionMetadata {
	return cidsdk.ActionMetadata{
		Name:        "applicationinspector-scan",
		Description: "Scans the repository for used features using Microsoft Application Inspector.",
		Category:    "sast",
		Scope:       cidsdk.ActionScopeProject,
		Rules:       []cidsdk.ActionRule{},
		Access: cidsdk.ActionAccess{
			Environment: []cidsdk.ActionAccessEnv{
				{
					Name:        "NCI_REPOSITORY_.*",
					Description: "The project properties to identify the repository.",
					Pattern:     true,
				},
				{
					Name:        "NCI_COMMIT_.*",
					Description: "The commit properties to identify the revision.",
					Pattern:     true,
				},
			},
			Executables: []cidsdk.ActionAccessExecutable{
				{
					Name: "appinspector",
				},
			},
		},
	}
}

func (a Action) Execute() (err error) {
	cfg := Config{}
	ctx, err := a.Sdk.ProjectAction(&cfg)
	if err != nil {
		return err
	}

	// files
	reportFile := cidsdk.JoinPath(ctx.Config.TempDir, "applicationinspector.json")

	// scan
	scanArgs := []string{
		"--no-show-progress",
		fmt.Sprintf("-s %q", ctx.ProjectDir),
		fmt.Sprintf("--base-path %q", ctx.ProjectDir),
		fmt.Sprintf("--repository-uri %q", ctx.Env["NCI_REPOSITORY_REMOTE"]),
		fmt.Sprintf("--commit-hash %q", ctx.Env["NCI_COMMIT_HASH"]),
		"-f json",
		fmt.Sprintf("-o %q", reportFile),
		"-g **/tests/**,**/.git/**,**/.dist/**,**/.tmp/**",
	}
	_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: `appinspector analyze ` + strings.Join(scanArgs, " "),
		WorkDir: ctx.ProjectDir,
	})
	if err != nil {
		return err
	}

	// store report
	err = a.Sdk.ArtifactUpload(cidsdk.ArtifactUploadRequest{
		File:          reportFile,
		Type:          "report",
		Format:        "applicationinspector",
		FormatVersion: "json",
	})
	if err != nil {
		return err
	}

	return nil
}
