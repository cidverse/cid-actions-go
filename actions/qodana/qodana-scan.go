package qodana

import (
	"fmt"
	"strings"

	cidsdk "github.com/cidverse/cid-sdk-go"
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
	if ctx.Module.BuildSystem == string(cidsdk.BuildSystemGradle) || ctx.Module.BuildSystem == string(cidsdk.BuildSystemMaven) {
		linter = "jvm"
	} else if ctx.Module.BuildSystem == string(cidsdk.BuildSystemNpm) {
		linter = "js"
	} else if ctx.Module.BuildSystem == string(cidsdk.BuildSystemGoMod) {
		linter = "go"
	} else if ctx.Module.BuildSystem == string(cidsdk.BuildSystemPipfile) || ctx.Module.BuildSystem == string(cidsdk.BuildSystemRequirementsTXT) || ctx.Module.BuildSystem == string(cidsdk.BuildSystemSetupPy) {
		linter = "python"
	} else if ctx.Module.BuildSystem == "dotnet" {
		linter = "dotnet"
	} else if ctx.Module.BuildSystem == "composer" {
		linter = "php"
	}
	if linter == "" {
		_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "warn", Message: "no supported linter, skipping!"})
		return nil
	}

	// qodana scan
	var scanOpts = []string{
		"--source-directory " + ctx.Module.ModuleDir,
		"--save-report",
		"--results-dir " + ctx.Config.TempDir,
	}
	scanResult, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		CaptureOutput: false,
		Command:       fmt.Sprintf("qodana-%s scan %s", linter, strings.Join(scanOpts, " ")),
		WorkDir:       ctx.ProjectDir,
		Env: map[string]string{
			"QODANA_TOKEN":      ctx.Env["QODANA_TOKEN"],
			"QODANA_REMOTE_URL": ctx.Env["NCI_REPOSITORY_REMOTE"],
			"QODANA_BRANCH":     ctx.Env["NCI_COMMIT_REF_NAME"],
			"QODANA_REVISION":   ctx.Env["NCI_COMMIT_SHA"],
			//"QODANA_JOB_URL":    ...,
		},
	})
	if err != nil {
		return err
	} else if scanResult.Code != 0 {
		return fmt.Errorf("qodana scan failed, exit code %d", scanResult.Code)
	}

	return nil
}
