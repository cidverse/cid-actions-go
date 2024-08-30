package ansiblelint

import (
	"path"

	cidsdk "github.com/cidverse/cid-sdk-go"
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
	_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: `ansible-lint .`,
		WorkDir: ctx.Module.ModuleDir,
	})
	if err != nil {
		return err
	}

	return nil
}
