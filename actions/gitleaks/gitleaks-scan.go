package gitleaks

import (
	"fmt"
	"path"
	"strings"

	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/owenrumney/go-sarif/v2/sarif"
)

type ScanAction struct {
	Sdk cidsdk.SDKClient
}

type ScanConfig struct {
}

func (a ScanAction) Execute() (err error) {
	cfg := ScanConfig{}
	ctx, err := a.Sdk.ProjectAction(&cfg)
	if err != nil {
		return err
	}

	// files
	reportFile := path.Join(ctx.Config.TempDir, "gitleaks.sarif.json")

	// opts
	var opts = []string{
		"gitleaks",
		"detect",
		"--source=.",
		"-v",
		"--no-git",
		"--report-format=sarif",
		"--report-path=" + reportFile,
		"--no-banner",
		"--exit-code 0",
	}
	if ctx.Env["CI"] == "true" {
		opts = append(opts, "--redact")
	}

	// scan
	scanResult, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: strings.Join(opts, " "),
		WorkDir: ctx.ProjectDir,
	})
	if err != nil {
		return err
	} else if scanResult.Code != 0 {
		return fmt.Errorf("gitleaks scan failed, exit code %d", scanResult.Code)
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
