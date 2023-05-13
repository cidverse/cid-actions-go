package helmlint

import (
	cidsdk "github.com/cidverse/cid-sdk-go"
)

type LintAction struct {
	Sdk cidsdk.SDKClient
}

type LintConfig struct {
}

func (a LintAction) Execute() (err error) {
	cfg := LintConfig{}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
	}

	if ctx.Module.BuildSystem == string(cidsdk.BuildSystemHelm) {
		_, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: `helm lint . --strict`,
			WorkDir: ctx.Module.ModuleDir,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
