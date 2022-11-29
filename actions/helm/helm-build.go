package helm

import (
	cidsdk "github.com/cidverse/cid-sdk-go"
)

type BuildAction struct {
	Sdk cidsdk.SDKClient
}

type BuildConfig struct {
}

func (a BuildAction) Execute() (err error) {
	chartDir := ""

	cfg := BuildConfig{}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
	}

	if ctx.Module.BuildSystem == string(cidsdk.BuildSystemHelm) {
		// restore the charts/ directory based on the Chart.lock file
		_, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: `helm dependency build .`,
			WorkDir: ctx.Module.ModuleDir,
		})
		if err != nil {
			return err
		}

		// package
		_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: `helm package ` + ctx.Module.ModuleDir + ` --version 0.0.1 --destination ` + chartDir,
			WorkDir: ctx.Module.ProjectDir,
		})
		if err != nil {
			return err
		}

		// update index
		_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: `helm repo index ` + chartDir,
			WorkDir: ctx.Module.ProjectDir,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
