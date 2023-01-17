package gitleaks

import (
	"os"
	"path"
	"strings"

	cidsdk "github.com/cidverse/cid-sdk-go"
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

	sarifDir := path.Join(ctx.Config.ArtifactDir, "gitleaks", "sarif")
	_ = os.MkdirAll(sarifDir, os.ModePerm)

	var opts []string
	opts = append(opts, "--source=.")
	opts = append(opts, "-v")
	opts = append(opts, "--no-git")
	opts = append(opts, "--report-format=sarif")
	opts = append(opts, "--report-path"+path.Join(sarifDir, "report.sarif"))
	if ctx.Env["CI"] == "true" {
		opts = append(opts, "--redact")
	}

	_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: strings.TrimRight(`gitleaks detect `+strings.Join(opts, " "), " "),
		WorkDir: ctx.ProjectDir,
	})
	if err != nil {
		return err
	}

	return nil
}
