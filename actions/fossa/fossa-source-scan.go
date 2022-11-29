package fossa

import (
	cidsdk "github.com/cidverse/cid-sdk-go"
)

type SourceScanAction struct {
	Sdk cidsdk.SDKClient
}

type SourceScanConfig struct {
}

func (a SourceScanAction) Execute() (err error) {
	cfg := SourceScanConfig{}
	ctx, err := a.Sdk.ProjectAction(&cfg)
	if err != nil {
		return err
	}

	_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: `fossa analyze`,
		WorkDir: ctx.ProjectDir,
	})
	if err != nil {
		return err
	}

	return nil
}
