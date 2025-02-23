package renovatelint

import (
	cidsdk "github.com/cidverse/cid-sdk-go"
)

type Action struct {
	Sdk cidsdk.SDKClient
}

type Config struct {
}

func (a Action) Metadata() cidsdk.ActionMetadata {
	return cidsdk.ActionMetadata{
		Name:        "renovate-lint",
		Description: "Lint the Renovate configuration file.",
		Category:    "sast",
		Scope:       cidsdk.ActionScopeProject,
		Rules: []cidsdk.ActionRule{
			{
				Type:       "cel",
				Expression: `contains(PROJECT_CONFIG_TYPES, "renovate")`,
			},
		},
		Access: cidsdk.ActionAccess{
			Environment: []cidsdk.ActionAccessEnv{},
			Executables: []cidsdk.ActionAccessExecutable{
				{
					Name: "renovate-config-validator",
				},
			},
		},
	}
}

func (a Action) Execute() (err error) {
	cfg := Config{}
	ctx, err := a.Sdk.ProjectAction(&cfg)
	if err != nil {
		return err
	}

	// run renovate-config-validator
	_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: `renovate-config-validator --strict .`,
		WorkDir: ctx.ProjectDir,
	})
	if err != nil {
		return err
	}

	return nil
}
