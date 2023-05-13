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
