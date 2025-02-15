package ggshield

import (
	"strings"

	cidsdk "github.com/cidverse/cid-sdk-go"
)

type Action struct {
	Sdk cidsdk.SDKClient
}

type ScanConfig struct {
}

const GitguardianPrefix = "GITGUARDIAN_"

func (a Action) Metadata() cidsdk.ActionMetadata {
	return cidsdk.ActionMetadata{
		Name:        "gitguardian-scan",
		Description: "Scans the repository for secrets using GitGuardian.",
		Category:    "sast",
		Scope:       cidsdk.ActionScopeProject,
		Rules: []cidsdk.ActionRule{
			{
				Type:       "cel",
				Expression: `ENV["GITGUARDIAN_API_KEY"] != "" && NCI_COMMIT_REF_TYPE == "branch"`,
			},
		},
		Access: cidsdk.ActionAccess{
			Environment: []cidsdk.ActionAccessEnv{
				{
					Name:    "GITGUARDIAN_.*",
					Pattern: true,
				},
			},
			Executables: []cidsdk.ActionAccessExecutable{
				{
					Name: "ggshield",
				},
			},
		},
	}
}

func (a Action) Execute() (err error) {
	cfg := ScanConfig{}
	ctx, err := a.Sdk.ProjectAction(&cfg)
	if err != nil {
		return err
	}

	// env
	execEnv := make(map[string]string)

	// GitGuardian env properties
	for key, value := range ctx.Env {
		if strings.HasPrefix(key, GitguardianPrefix) {
			execEnv[key] = value
		}
	}
	_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: `ggshield scan path -r -y .`,
		WorkDir: ctx.ProjectDir,
		Env:     execEnv,
	})
	if err != nil {
		return err
	}

	return nil
}
