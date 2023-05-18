package scorecardscan

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

func (a Action) Execute() (err error) {
	cfg := Config{}
	ctx, err := a.Sdk.ProjectAction(&cfg)
	if err != nil {
		return err
	}

	// TODO: remove
	if _, ok := ctx.Env["NCI_REPOSITORY_URL"]; !ok {
		ctx.Env["NCI_REPOSITORY_URL"] = strings.TrimSuffix(ctx.Env["NCI_REPOSITORY_REMOTE"], ".git")
	}

	// scorecard scan
	var scanOpts = []string{
		fmt.Sprintf(`--repo %q`, ctx.Env["NCI_REPOSITORY_URL"]),
		`--format json`,
		fmt.Sprintf(`--commit %q`, ctx.Env["NCI_COMMIT_HASH"]),
		`--checks "Contributors,Dependency-Update-Tool,Maintained,Security-Policy,Fuzzing,Branch-Protection,CI-Tests,Signed-Releases,Binary-Artifacts,SAST,License,Pinned-Dependencies,CII-Best-Practices,Code-Review,Dangerous-Workflow,Packaging,Token-Permissions,Vulnerabilities"`,
	}
	scanEnv := map[string]string{}
	if ctx.Env["NCI_REPOSITORY_HOST_TYPE"] == "github" {
		scanEnv["GITHUB_AUTH_TOKEN"] = ctx.Env["GITHUB_TOKEN"]
	} else if ctx.Env["NCI_REPOSITORY_HOST_TYPE"] == "gitlab" {
		scanEnv["GITLAB_AUTH_TOKEN"] = ctx.Env["GITLAB_TOKEN"]
	}
	scanResult, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		CaptureOutput: true,
		Command:       fmt.Sprintf("scorecard %s", strings.Join(scanOpts, " ")),
		WorkDir:       ctx.ProjectDir,
		Env:           scanEnv,
	})
	if err != nil {
		return err
	} else if scanResult.Code != 0 {
		return fmt.Errorf("qodana scan failed, exit code %d", scanResult.Code)
	}

	// parse / validate report
	report, err := sarif.FromBytes([]byte(scanResult.Stdout))
	if err != nil {
		return err
	}

	// store result
	err = a.Sdk.ArtifactUpload(cidsdk.ArtifactUploadRequest{
		File:          "ossf-scorecard.sarif.json",
		Content:       scanResult.Stdout,
		Type:          "report",
		Format:        "sarif",
		FormatVersion: report.Version,
	})
	if err != nil {
		return err
	}

	return nil
}
