package gitleaksscan

import (
	"fmt"
	"strings"

	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/owenrumney/go-sarif/v3/pkg/report/v210/sarif"
)

type Action struct {
	Sdk cidsdk.SDKClient
}

type Config struct {
}

func (a Action) Metadata() cidsdk.ActionMetadata {
	return cidsdk.ActionMetadata{
		Name:        "gitleaks-scan",
		Description: "Scans the repository for secrets using Gitleaks.",
		Category:    "sast",
		Scope:       cidsdk.ActionScopeProject,
		Rules: []cidsdk.ActionRule{
			{
				Type:       "cel",
				Expression: `NCI_COMMIT_REF_TYPE == "branch" && size(PROJECT_BUILD_SYSTEMS) > 0`,
			},
		},
		Access: cidsdk.ActionAccess{
			Environment: []cidsdk.ActionAccessEnv{},
			Executables: []cidsdk.ActionAccessExecutable{
				{
					Name: "gitleaks",
				},
				{
					Name: "gitlab-sarif-converter",
				},
			},
		},
		Output: cidsdk.ActionOutput{
			Artifacts: []cidsdk.ActionArtifactType{
				{
					Type:   "report",
					Format: "sarif",
				},
				{
					Type:   "report",
					Format: "gl-codequality",
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
	reportFile := cidsdk.JoinPath(ctx.Config.TempDir, "gitleaks.sarif.json")

	// opts
	var opts = []string{
		"gitleaks", "detect",
		"--source=.",
		"-v",
		"--no-git",
		"--report-format=sarif",
		fmt.Sprintf("--report-path=%q", reportFile),
		"--no-banner",
		"--redact=85", // redact 85% of the secret
		"--exit-code 0",
	}

	// scan
	cmdResult, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: strings.Join(opts, " "),
		WorkDir: ctx.ProjectDir,
	})
	if err != nil {
		return err
	} else if cmdResult.Code != 0 {
		return fmt.Errorf("gitleaks scan failed, exit code %d", cmdResult.Code)
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

	// gitlab conversion
	if ctx.Env["NCI_REPOSITORY_HOST_TYPE"] == "gitlab" {
		// code-quality report
		codeQualityFile := cidsdk.JoinPath(ctx.Config.TempDir, "gl-code-quality-report.json")
		cmdResult, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: fmt.Sprintf("gitlab-sarif-converter --type=codequality %q %q", reportFile, codeQualityFile),
			WorkDir: ctx.ProjectDir,
		})
		if err != nil {
			return err
		} else if cmdResult.Code != 0 {
			return fmt.Errorf("gitlab-sarif-converter failed, exit code %d", cmdResult.Code)
		}

		err = a.Sdk.ArtifactUpload(cidsdk.ArtifactUploadRequest{
			File:   codeQualityFile,
			Type:   "report",
			Format: "gl-codequality",
		})
		if err != nil {
			return err
		}
	}

	return nil
}
