package pythonlint

import (
	cidsdk "github.com/cidverse/cid-sdk-go"
)

type Action struct {
	Sdk cidsdk.SDKClient
}

type LintConfig struct {
}

func (a Action) Metadata() cidsdk.ActionMetadata {
	return cidsdk.ActionMetadata{
		Name:        "python-lint",
		Description: "Runs the flake8 tool to lint your python project.",
		Category:    "lint",
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
					Name: "flake8",
				},
			},
		},
	}
}

func (a Action) Execute() (err error) {
	cfg := LintConfig{}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
	}

	// TODO: fix
	_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: `flake8 .`,
		WorkDir: ctx.Module.ModuleDir,
	})
	if err != nil {
		return err
	}

	return nil
}
