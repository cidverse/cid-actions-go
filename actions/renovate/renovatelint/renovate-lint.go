package renovatelint

import (
	cidsdk "github.com/cidverse/cid-sdk-go"
)

type Action struct {
	Sdk cidsdk.SDKClient
}

type Config struct {
}

func (a Action) Execute() (err error) {
	cfg := Config{}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
	}

	// run renovate-config-validator
	_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: `renovate-config-validator --strict`,
		WorkDir: ctx.Module.ModuleDir,
	})
	if err != nil {
		return err
	}

	return nil
}
