package renovatelint

import (
	"fmt"

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
	// query action data
	d, err := a.Sdk.ProjectActionDataV1()
	if err != nil {
		return err
	}

	// parse config
	cfg := Config{}
	cidsdk.PopulateFromEnv(&cfg, d.Env)

	// run renovate-config-validator
	cmdResult, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: `renovate-config-validator --strict .`,
		WorkDir: d.ProjectDir,
	})
	if err != nil {
		return err
	} else if cmdResult.Code != 0 {
		return fmt.Errorf("command failed, exit code %d", cmdResult.Code)
	}

	return nil
}
