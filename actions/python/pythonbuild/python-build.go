package pythonbuild

import (
	cidsdk "github.com/cidverse/cid-sdk-go"
)

type BuildAction struct {
	Sdk cidsdk.SDKClient
}

type BuildConfig struct {
}

func (a BuildAction) Execute() (err error) {
	cfg := BuildConfig{}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
	}

	if ctx.Module.BuildSystem == string(cidsdk.BuildSystemRequirementsTXT) {
		_, installErr := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: `pip install -r requirements.txt`,
			WorkDir: ctx.Module.ModuleDir,
		})
		if installErr != nil {
			return installErr
		}
	} else if ctx.Module.BuildSystem == string(cidsdk.BuildSystemPipfile) {
		_, installErr := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: `pipenv install`,
			WorkDir: ctx.Module.ModuleDir,
		})
		if installErr != nil {
			return installErr
		}
	} else if ctx.Module.BuildSystem == string(cidsdk.BuildSystemSetupPy) {
		_, installErr := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: `pip install .`,
			WorkDir: ctx.Module.ModuleDir,
		})
		if installErr != nil {
			return installErr
		}
	}

	return nil
}
