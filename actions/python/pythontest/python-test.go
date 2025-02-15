package pythontest

import (
	cidsdk "github.com/cidverse/cid-sdk-go"
)

type Action struct {
	Sdk cidsdk.SDKClient
}

type TestConfig struct {
}

func (a Action) Metadata() cidsdk.ActionMetadata {
	return cidsdk.ActionMetadata{
		Name:        "python-test",
		Description: "Runs all tests in your python project.",
		Category:    "test",
		Scope:       cidsdk.ActionScopeModule,
		Rules: []cidsdk.ActionRule{
			{
				Type:       "cel",
				Expression: `MODULE_BUILD_SYSTEM == "python-requirements.txt"`,
			},
			{
				Type:       "cel",
				Expression: `MODULE_BUILD_SYSTEM == "pipfile"`,
			},
			{
				Type:       "cel",
				Expression: `MODULE_BUILD_SYSTEM == "setup.py"`,
			},
		},
		Access: cidsdk.ActionAccess{
			Environment: []cidsdk.ActionAccessEnv{},
			Executables: []cidsdk.ActionAccessExecutable{
				{
					Name: "pytest",
				},
			},
		},
	}
}

func (a Action) Execute() (err error) {
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
