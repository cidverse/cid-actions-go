package semgrepscan

import (
	"fmt"
	"strconv"
	"strings"

	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/go-playground/validator/v10"
)

type Action struct {
	Sdk cidsdk.SDKClient
}

type Config struct {
	RuleSets []string
}

func (a Action) Metadata() cidsdk.ActionMetadata {
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
					Name:        "SEMGREP_RULES",
					Description: "See option --config.",
				},
				{
					Name:        "SEMGREP_APP_TOKEN",
					Description: "Semgrep AppSec Platform Token",
					Secret:      true,
				},
			},
			Executables: []cidsdk.ActionAccessExecutable{
				{
					Name: "semgrep",
				},
				{
					Name: "gitlab-sarif-converter",
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
				{
					Type:   "report",
					Format: "gl-codequality",
				},
			},
		},
	}
}

func (a Action) GetConfig(d *cidsdk.ProjectActionData) (Config, error) {
	cfg := Config{}
	cidsdk.PopulateFromEnv(&cfg, d.Env)

	// validate
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func (a Action) Execute() (err error) {
	// query action data
	d, err := a.Sdk.ProjectActionDataV1()
	if err != nil {
		return err
	}

	// parse config
	cfg, err := a.GetConfig(d)
	if err != nil {
		return err
	}

	// files
	reportFile := cidsdk.JoinPath(d.Config.TempDir, "semgrep.sarif.json")

	// defaults
	if len(cfg.RuleSets) == 0 {
		cfg.RuleSets = []string{"p/ci"}
	}

	// scan
	var opts = []string{
		"semgrep", "ci",
		"--text", // output plain text format in stdout
		"--sarif-output=" + strconv.Quote(reportFile), // output sarif format to file
		"--metrics=off",
		"--disable-version-check",
		"--exclude=.dist",
		"--exclude=.tmp",
	}
	/*
		if val, ok := ctx.Env["NCI_MERGE_REQUEST_SOURCE_HASH"]; ok && len(val) > 0 {
			opts = append(opts, "--baseline", val)
		}
	*/

	// ruleSets
	for _, config := range cfg.RuleSets {
		opts = append(opts, "--config", strconv.Quote(config))
	}

	// scan
	cmdResult, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: strings.Join(opts, " "),
		WorkDir: d.ProjectDir,
		Env: map[string]string{
			"SEMGREP_RULES":     d.Env["SEMGREP_RULES"],
			"SEMGREP_APP_TOKEN": d.Env["SEMGREP_APP_TOKEN"],
		},
	})
	if err != nil {
		return err
	} else if cmdResult.Code != 0 {
		return fmt.Errorf("failed, exit code %d. error: %s", cmdResult.Code, cmdResult.Stderr)
	}

	// store report
	err = a.Sdk.ArtifactUpload(cidsdk.ArtifactUploadRequest{
		File:          reportFile,
		Type:          "report",
		Format:        "sarif",
		FormatVersion: "2.1.0",
	})
	if err != nil {
		return fmt.Errorf("failed to upload report %s: %w", reportFile, err)
	}

	// gitlab conversion
	if d.Env["NCI_REPOSITORY_HOST_TYPE"] == "gitlab" {
		// code-quality report
		codeQualityFile := cidsdk.JoinPath(d.Config.TempDir, "gl-code-quality-report.json")
		cmdResult, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: fmt.Sprintf("gitlab-sarif-converter --type=codequality %q %q", reportFile, codeQualityFile),
			WorkDir: d.ProjectDir,
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
