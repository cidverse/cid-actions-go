package semgrepscan

import (
	"fmt"
	"strings"

	cidsdk "github.com/cidverse/cid-sdk-go"
)

type ScanAction struct {
	Sdk cidsdk.SDKClient
}

type ScanConfig struct {
	RuleSets []string
}

func (a ScanAction) Metadata() cidsdk.ActionMetadata {
	return cidsdk.ActionMetadata{
		Name:        "semgrep-scan",
		Description: "Scans the repository for security issues using semgrep.",
		Category:    "sast",
		Scope:       cidsdk.ActionScopeProject,
		Rules: []cidsdk.ActionRule{
			{
				Type:       "cel",
				Expression: `NCI_COMMIT_REF_TYPE == "branch" && size(PROJECT_BUILD_SYSTEMS) > 0`,
			},
		},
		Access: cidsdk.ActionAccess{
			Environment: []cidsdk.ActionAccessEnv{
				{
					Name:        "SEMGREP_.*",
					Description: "Semgrep configuration properties",
					Pattern:     true,
				},
			},
			Executables: []cidsdk.ActionAccessExecutable{
				{
					Name: "semgrep",
				},
			},
			Network: []cidsdk.ActionAccessNetwork{
				{
					Host: "semgrep.dev:443",
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

func (a ScanAction) Execute() (err error) {
	cfg := ScanConfig{}
	ctx, err := a.Sdk.ProjectAction(&cfg)
	if err != nil {
		return err
	}

	// files
	reportFile := cidsdk.JoinPath(ctx.Config.TempDir, "semgrep.sarif.json")

	// defaults
	if len(cfg.RuleSets) == 0 {
		cfg.RuleSets = []string{"p/ci"}
	}

	// scan
	var opts = []string{"semgrep", "ci", "--sarif", "--quiet", "--metrics=off", "--disable-version-check", "--exclude=.dist", "--exclude=.tmp"}
	/*
		if val, ok := ctx.Env["NCI_MERGE_REQUEST_SOURCE_HASH"]; ok && len(val) > 0 {
			opts = append(opts, "--baseline", val)
		}
	*/

	// ruleSets
	for _, config := range cfg.RuleSets {
		opts = append(opts, fmt.Sprintf("--config %q", config))
	}

	// scan
	scanResult, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command:       strings.Join(opts, " "),
		WorkDir:       ctx.ProjectDir,
		CaptureOutput: true,
	})
	if err != nil {
		return err
	} else if scanResult.Code != 0 {
		return fmt.Errorf("failed, exit code %d. error: %s", scanResult.Code, scanResult.Stderr)
	}

	_ = a.Sdk.FileWrite(reportFile, []byte(scanResult.Stdout))

	// store report
	err = a.Sdk.ArtifactUpload(cidsdk.ArtifactUploadRequest{
		File:          reportFile,
		Type:          "report",
		Format:        "sarif",
		FormatVersion: "2.1.0",
	})
	if err != nil {
		return err
	}

	return nil
}
