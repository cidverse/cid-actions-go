package ggshield

import (
	cidsdk "github.com/cidverse/cid-sdk-go"
)

type ScanAction struct {
	Sdk cidsdk.SDKClient
}

type ScanConfig struct {
}

const GitguardianAPIKey = "GITGUARDIAN_API_KEY"
const GitguardianPrefix = "GITGUARDIAN_"

func (a ScanAction) Execute() (err error) {
	cfg := ScanConfig{}
	ctx, err := a.Sdk.ProjectAction(&cfg)
	if err != nil {
		return err
	}

	/*
	   // env
	   	execEnv := ctx.Env
	   	execEnv[GitguardianAPIKey] = api.GetEnvValue(ctx, GitguardianAPIKey)

	   	// GitGuardian env properties
	   	for key, value := range ctx.Env {
	   		if strings.HasPrefix(key, GitguardianPrefix) {
	   			execEnv[key] = value
	   		}
	   	}
	*/

	_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: `ggshield scan path -r -y .`,
		WorkDir: ctx.ProjectDir,
	})
	if err != nil {
		return err
	}

	return nil
}
