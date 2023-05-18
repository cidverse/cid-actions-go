package qodanascan

import (
	"fmt"
	"strings"

	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/owenrumney/go-sarif/v2/sarif"
	"github.com/thoas/go-funk"
)

type ScanAction struct {
	Sdk cidsdk.SDKClient
}

type ScanConfig struct {
	QodanaToken string `json:"qodana_token"  env:"QODANA_TOKEN"`
}

func (a ScanAction) Execute() (err error) {
	cfg := ScanConfig{}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
	}

	// choose linter
	linter := ""
	if ctx.Module.Language != nil {
		if funk.Contains(*ctx.Module.Language, string(cidsdk.LanguageJava)) {
			linter = "jvm"
		} else if funk.Contains(*ctx.Module.Language, string(cidsdk.LanguageJavascript)) || funk.Contains(*ctx.Module.Language, string(cidsdk.LanguageTypescript)) {
			linter = "js"
		} else if funk.Contains(*ctx.Module.Language, string(cidsdk.LanguageGolang)) {
			linter = "go"
		} else if funk.Contains(*ctx.Module.Language, string(cidsdk.LanguagePython)) {
			linter = "python"
		} else if funk.Contains(*ctx.Module.Language, string(cidsdk.LanguagePHP)) {
			linter = "php"
		}
	}
	if ctx.Module.BuildSystem == string(cidsdk.BuildSystemDotNet) {
		linter = "dotnet"
	}
	if linter == "" {
		_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "warn", Message: "no supported linter, skipping!"})
		return nil
	}
	_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "info", Message: "using qodana linter", Context: map[string]interface{}{"linter": "qodana-" + linter}})

	// qodana scan
	var scanOpts = []string{
		"--source-directory=" + ctx.Module.ModuleDir,
		"--results-dir=" + ctx.Config.TempDir,
		"--fail-threshold 10000",
	}
	scanResult, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		CaptureOutput: false,
		Command:       fmt.Sprintf("qodana-%s %s", linter, strings.Join(scanOpts, " ")),
		WorkDir:       ctx.Module.ModuleDir,
		Env: map[string]string{
			"QODANA_TOKEN":      ctx.Env["QODANA_TOKEN"],
			"QODANA_REMOTE_URL": ctx.Env["NCI_REPOSITORY_REMOTE"],
			"QODANA_BRANCH":     ctx.Env["NCI_COMMIT_REF_NAME"],
			"QODANA_REVISION":   ctx.Env["NCI_COMMIT_HASH"],
			//"QODANA_JOB_URL":    ...,
		},
	})
	if err != nil {
		return err
	} else if scanResult.Code != 0 {
		return fmt.Errorf("qodana scan failed, exit code %d", scanResult.Code)
	}

	// parse / validate report
	qodanaReportFile := fmt.Sprintf("%s/qodana.sarif.json", ctx.Config.TempDir)
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
		File:          fmt.Sprintf("%s/qodana.sarif.json", ctx.Config.TempDir),
		Module:        ctx.Module.Slug,
		Type:          "report",
		Format:        "sarif",
		FormatVersion: report.Version,
	})
	if err != nil {
		return err
	}

	return nil
}
