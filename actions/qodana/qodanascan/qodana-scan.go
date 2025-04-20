package qodanascan

import (
	"fmt"
	"strings"

	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/go-playground/validator/v10"
	"github.com/owenrumney/go-sarif/v3/pkg/report/v210/sarif"
)

type Action struct {
	Sdk cidsdk.SDKClient
}

type Config struct {
	QodanaToken              string `json:"qodana_token"                 env:"QODANA_TOKEN"`
	QodanaUltimate           bool   `json:"qodana_ultimate"              env:"QODANA_ULTIMATE"`
	QodanaEarlyAccessPreview bool   `json:"qodana_early_access_preview"  env:"QODANA_EAP"`
}

func (a Action) Metadata() cidsdk.ActionMetadata {
	return cidsdk.ActionMetadata{
		Name:          "qodana-scan",
		Description:   "Scans the repository for security issues using JetBrains Qodana.",
		Documentation: ``,
		Category:      "sast",
		Scope:         cidsdk.ActionScopeModule,
		Rules: []cidsdk.ActionRule{
			{
				Type:       "cel",
				Expression: `ENV["QODANA_TOKEN"] != "" && MODULE_BUILD_SYSTEM != "" && NCI_COMMIT_REF_TYPE == "branch"`,
			},
		},
		Access: cidsdk.ActionAccess{
			Environment: []cidsdk.ActionAccessEnv{
				{
					Name:        "QODANA_TOKEN",
					Description: "The Qodana cloud project token.",
					Required:    true,
					Secret:      true,
				},
				{
					Name:        "QODANA_ULTIMATE",
					Description: "Set if you have a Qodana Ultimate license, will use the commercial IDEs instead of the community edition.",
				},
				{
					Name:        "QODANA_EAP",
					Description: "Enable Qodana Early Access Preview IDEs, does not require a license.",
				},
				{
					Name:        "NCI_REPOSITORY_.*",
					Description: "The project properties sonar needs to identify the repository.",
					Pattern:     true,
				},
				{
					Name:        "NCI_COMMIT_.*",
					Description: "The commit properties sonar needs to identify the revision.",
					Pattern:     true,
				},
			},
			Executables: []cidsdk.ActionAccessExecutable{
				{
					Name: "qodana",
				},
				{
					Name: "gitlab-sarif-converter",
				},
			},
		},
	}
}

func (a Action) GetConfig(d *cidsdk.ModuleActionData) (Config, error) {
	cfg := Config{}
	cidsdk.PopulateFromEnv(&cfg, d.Env)

	// enable eap by default
	cfg.QodanaEarlyAccessPreview = true

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
	d, err := a.Sdk.ModuleActionDataV1()
	if err != nil {
		return err
	}

	// parse config
	cfg, err := a.GetConfig(d)
	if err != nil {
		return err
	}

	// choose ide
	ideName := ""
	if d.Module.Language != nil {
		if (*d.Module.Language)[string(cidsdk.LanguagePython)] != "" {
			if cfg.QodanaUltimate {
				ideName = "QDPY"
			} else {
				ideName = "QDPYC"
			}
		} else if (*d.Module.Language)[string(cidsdk.LanguagePHP)] != "" {
			if cfg.QodanaEarlyAccessPreview {
				ideName = "QDPHP-EAP"
			} else if cfg.QodanaUltimate {
				ideName = "QDPHP"
			}
		}
	}
	if d.Module.BuildSystem == string(cidsdk.BuildSystemDotNet) {
		if cfg.QodanaEarlyAccessPreview {
			ideName = "QDNET-EAP"
		} else if cfg.QodanaUltimate {
			ideName = "QDNET"
		}
	} else if d.Module.BuildSystem == string(cidsdk.BuildSystemGoMod) {
		if cfg.QodanaEarlyAccessPreview {
			ideName = "QDGO-EAP"
		} else if cfg.QodanaUltimate {
			ideName = "QDGO"
		}
	} else if d.Module.BuildSystem == string(cidsdk.BuildSystemMaven) || d.Module.BuildSystem == string(cidsdk.BuildSystemGradle) {
		if cfg.QodanaEarlyAccessPreview {
			ideName = "QDJVM-EAP"
		} else if cfg.QodanaUltimate {
			ideName = "QDJVM"
		} else {
			ideName = "QDJVMC"
		}
	} else if d.Module.BuildSystem == string(cidsdk.BuildSystemNpm) {
		if cfg.QodanaEarlyAccessPreview {
			ideName = "QDJS-EAP"
		} else if cfg.QodanaUltimate {
			ideName = "QDJS"
		}
	}

	if ideName == "" {
		_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "info", Message: "no supported linter, skipping!"})
		return nil
	}
	_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "info", Message: "using qodana linter", Context: map[string]interface{}{"qodana-type": ideName}})

	// qodana scan
	var scanOpts = []string{
		"--ide=" + ideName,
		"--project-dir=" + d.ProjectDir,
		"--source-directory=" + d.Module.ModuleDir,
		"--results-dir=" + d.Config.TempDir,
		"--fail-threshold 10000",
	}
	scanResult, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		CaptureOutput: false,
		Command:       fmt.Sprintf("qodana scan %s", strings.Join(scanOpts, " ")),
		Constraint:    "",
		WorkDir:       d.Module.ModuleDir,
		Env: map[string]string{
			"QODANA_TOKEN":      d.Env["QODANA_TOKEN"],
			"QODANA_REMOTE_URL": d.Env["NCI_REPOSITORY_REMOTE"],
			"QODANA_BRANCH":     d.Env["NCI_COMMIT_REF_NAME"],
			"QODANA_REVISION":   d.Env["NCI_COMMIT_HASH"],
			//"QODANA_JOB_URL":    ...,
		},
	})
	if err != nil {
		return err
	} else if scanResult.Code != 0 {
		return fmt.Errorf("qodana scan failed, exit code %d", scanResult.Code)
	}

	// parse / validate report
	qodanaReportFile := fmt.Sprintf("%s/qodana.sarif.json", d.Config.TempDir)
	content, err := a.Sdk.FileRead(qodanaReportFile)
	if err != nil {
		return err
	}
	report, err := sarif.FromBytes([]byte(content))
	if err != nil {
		return err
	}

	// store result
	err = a.Sdk.ArtifactUpload(cidsdk.ArtifactUploadRequest{
		File:          fmt.Sprintf("%s/qodana.sarif.json", d.Config.TempDir),
		Module:        d.Module.Slug,
		Type:          "report",
		Format:        "sarif",
		FormatVersion: report.Version,
	})
	if err != nil {
		return err
	}

	return nil
}
