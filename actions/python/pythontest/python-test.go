package pythontest

import (
	cidsdk "github.com/cidverse/cid-sdk-go"
)

type TestAction struct {
	Sdk cidsdk.SDKClient
}

type TestConfig struct {
}

func (a TestAction) Execute() (err error) {
	cfg := TestConfig{}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
	}

	// https://docs.pytest.org/en/6.2.x/
	if ctx.Module.BuildSystem == string(cidsdk.BuildSystemRequirementsTXT) || ctx.Module.BuildSystem == string(cidsdk.BuildSystemPipfile) || ctx.Module.BuildSystem == string(cidsdk.BuildSystemSetupPy) {
		_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: `pytest`,
			WorkDir: ctx.Module.ModuleDir,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
