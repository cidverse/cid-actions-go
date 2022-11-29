package java

import (
	cidsdk "github.com/cidverse/cid-sdk-go"
)

type BuildAction struct {
	Sdk cidsdk.SDKClient
}

type BuildConfig struct {
}

func (a BuildAction) Execute() (err error) {
	cfg := BuildConfig{}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
	}

	if ctx.Module.BuildSystem == string(cidsdk.BuildSystemGradle) {

	} else if ctx.Module.BuildSystem == string(cidsdk.BuildSystemMaven) {

	}

	return nil
}
