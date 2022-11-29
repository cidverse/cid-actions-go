package java

import (
	cidsdk "github.com/cidverse/cid-sdk-go"
)

type TestAction struct {
	Sdk cidsdk.SDKClient
}

type TestConfig struct {
}

func (a TestAction) Execute() (err error) {
	cfg := TestConfig{}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
	}

	if ctx.Module.BuildSystem == string(cidsdk.BuildSystemGradle) {

	} else if ctx.Module.BuildSystem == string(cidsdk.BuildSystemMaven) {

	}

	return nil
}
