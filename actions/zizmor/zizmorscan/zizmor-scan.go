package zizmorscan

import (
	"fmt"
	"strings"

	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/owenrumney/go-sarif/v2/sarif"
)

type Action struct {
	Sdk cidsdk.SDKClient
}

type Config struct {
}

func (a Action) Metadata() cidsdk.ActionMetadata {
	return cidsdk.ActionMetadata{
		Name:        "zizmor-scan",
		Description: "A static analysis tool for GitHub Actions",
		Category:    "sast",
		Scope:       cidsdk.ActionScopeProject,
		Links: map[string]string{
			"repo": "https://github.com/woodruffw/zizmor",
			"docs": "https://woodruffw.github.io/zizmor/",
		},
		Rules: []cidsdk.ActionRule{
			{
				Type:       "cel",
				Expression: `contains(PROJECT_CONFIG_TYPES, "github-workflow")`,
			},
		},
		Access: cidsdk.ActionAccess{
			Environment: []cidsdk.ActionAccessEnv{
				{
					Name:        "GH_HOSTNAME",
					Description: "GH_HOSTNAME is required for some online audits",
				},
				{
					Name:        "GH_TOKEN",
					Description: "GH_TOKEN is required for some online audits",
				},
			},
			Executables: []cidsdk.ActionAccessExecutable{
				{
					Name: "zizmor",
				},
			},
		},
		Output: cidsdk.ActionOutput{
			Artifacts: []cidsdk.ActionArtifactType{
				{
					Type:   "report",
					Format: "sarif",
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
	reportFile := cidsdk.JoinPath(ctx.Config.TempDir, "zizmor.sarif.json")

	// scan
	var opts = []string{".", "--format", "sarif", "-q"}
	resp, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command:       strings.TrimRight(`zizmor `+strings.Join(opts, " "), " "),
		WorkDir:       ctx.ProjectDir,
		CaptureOutput: true,
	})
	if err != nil {
		return err
	}

	// write and parse report
	sarifContent := []byte(resp.Stdout)
	err = a.Sdk.FileWrite(reportFile, sarifContent)
	if err != nil {
		return fmt.Errorf("failed to write report content to file %s: %s", reportFile, err.Error())
	}
	report, err := sarif.FromBytes(sarifContent)
	if err != nil {
		return err
	}

	// store report
	err = a.Sdk.ArtifactUpload(cidsdk.ArtifactUploadRequest{
		File:          reportFile,
		Type:          "report",
		Format:        "sarif",
		FormatVersion: report.Version,
	})
	if err != nil {
		return err
	}

	return nil
}
