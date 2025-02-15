package pythonbuild

import (
	cidsdk "github.com/cidverse/cid-sdk-go"
)

type Action struct {
	Sdk cidsdk.SDKClient
}

type BuildConfig struct {
}

func (a Action) Metadata() cidsdk.ActionMetadata {
	return cidsdk.ActionMetadata{
		Name:        "python-build",
		Description: "Builds the python project.",
		Category:    "build",
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
			Environment: []cidsdk.ActionAccessEnv{
				{
					Name:        "SEMGREP_.*",
					Description: "Semgrep configuration properties",
					Pattern:     true,
				},
			},
			Executables: []cidsdk.ActionAccessExecutable{
				{
					Name: "pip",
				},
				{
					Name: "pipenv",
				},
			},
		},
	}
}

func (a Action) Execute() (err error) {
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
