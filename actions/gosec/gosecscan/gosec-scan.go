package gosecscan

import (
	"fmt"
	"strings"

	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/owenrumney/go-sarif/v2/sarif"
)

type Action struct {
	Sdk cidsdk.SDKClient
}

type ScanConfig struct {
}

func (a Action) Metadata() cidsdk.ActionMetadata {
	return cidsdk.ActionMetadata{
		Name:        "gosec-scan",
		Description: `Scans the repository for security issues using gosec.`,
		Category:    "sast",
		Scope:       cidsdk.ActionScopeModule,
		Rules: []cidsdk.ActionRule{
			{
				Type:       "cel",
				Expression: `MODULE_BUILD_SYSTEM == "gomod"`,
			},
		},
		Access: cidsdk.ActionAccess{
			Environment: []cidsdk.ActionAccessEnv{},
			Executables: []cidsdk.ActionAccessExecutable{
				{
					Name: "gosec",
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
	cfg := ScanConfig{}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
	}

	// files
	reportFile := cidsdk.JoinPath(ctx.Config.TempDir, "gosec.sarif.json")

	// scan
	var opts = []string{"-no-fail", "-fmt sarif", "-out " + reportFile, "./..."}
	_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: strings.TrimRight(`gosec `+strings.Join(opts, " "), " "),
		WorkDir: ctx.ProjectDir,
	})
	if err != nil {
		return err
	}

	// parse report
	reportContent, err := a.Sdk.FileRead(reportFile)
	if err != nil {
		return fmt.Errorf("failed to read report content from file %s: %s", reportFile, err.Error())
	}
	report, err := sarif.FromBytes([]byte(reportContent))
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
