package ansiblelint

import (
	"fmt"
	"path"

	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/owenrumney/go-sarif/v2/sarif"
)

type Action struct {
	Sdk cidsdk.SDKClient
}

type ScanConfig struct {
}

func (a Action) Execute() (err error) {
	cfg := ScanConfig{}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
	}

	// files
	reportFile := cidsdk.JoinPath(ctx.Config.TempDir, "ansiblelint.sarif.json")

	// role and collection requirements
	if a.Sdk.FileExists(path.Join(ctx.Module.ModuleDir, "requirements.yml")) {
		_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: `ansible-galaxy collection install -r requirements.yml`,
			WorkDir: ctx.Module.ModuleDir,
		})
		if err != nil {
			return err
		}
	}

	// lint
	// config lookup: https://ansible.readthedocs.io/projects/lint/configuring/#using-local-configuration-files
	_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: fmt.Sprintf(`ansible-lint --project . --sarif-file %s`, reportFile),
		WorkDir: ctx.Module.ModuleDir,
	})
	if err != nil {
		return err
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
