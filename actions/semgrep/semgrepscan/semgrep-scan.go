package semgrepscan

import (
	"fmt"
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
	reportFile := cidsdk.JoinPath(ctx.Config.TempDir, "semgrep.sarif.json")

	// scan
	var opts = []string{"semgrep", "scan", "--config p/ci", "--sarif", "--quiet", "--metrics=off", "--disable-version-check", "--exclude=.dist", "--exclude=.tmp"}
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

	// parse report
	report, err := sarif.FromBytes([]byte(scanResult.Stdout))
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
