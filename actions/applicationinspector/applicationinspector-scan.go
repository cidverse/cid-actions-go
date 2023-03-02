package applicationinspector

import (
	"fmt"
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

	// files
	reportFile := cidsdk.JoinPath(ctx.Config.TempDir, "applicationinspector.json")

	// scan
	scanArgs := []string{
		"--no-show-progress",
		fmt.Sprintf("-s %q", ctx.ProjectDir),
		fmt.Sprintf("--base-path %q", ctx.ProjectDir),
		fmt.Sprintf("--repository-uri %q", ctx.Env["NCI_REPOSITORY_REMOTE"]),
		fmt.Sprintf("--commit-hash %q", ctx.Env["NCI_COMMIT_SHA"]),
		"-f json",
		fmt.Sprintf("-o %q", reportFile),
		"-g **/tests/**,**/.git/**,**/.dist/**,**/.tmp/**",
	}
	_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: `appinspector analyze ` + strings.Join(scanArgs, " "),
		WorkDir: ctx.ProjectDir,
	})
	if err != nil {
		return err
	}

	// store report
	err = a.Sdk.ArtifactUpload(cidsdk.ArtifactUploadRequest{
		File:          reportFile,
		Type:          "report",
		Format:        "applicationinspector",
		FormatVersion: "json",
	})
	if err != nil {
		return err
	}

	return nil
}
