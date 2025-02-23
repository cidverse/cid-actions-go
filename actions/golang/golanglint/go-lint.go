package golanglint

import (
	"errors"

	cidsdk "github.com/cidverse/cid-sdk-go"
)

type Action struct {
	Sdk cidsdk.SDKClient
}

func (a Action) Metadata() cidsdk.ActionMetadata {
	return cidsdk.ActionMetadata{
		Name:        "go-lint",
		Description: "Runs the golangci-lint tool on your go project.",
		Category:    "sast",
		Scope:       cidsdk.ActionScopeModule,
		Rules: []cidsdk.ActionRule{
			{
				Type:       "cel",
				Expression: `MODULE_BUILD_SYSTEM == "gomod"`,
			},
		},
		Access: cidsdk.ActionAccess{
			Environment: []cidsdk.ActionAccessEnv{},
			Executables: []cidsdk.ActionAccessExecutable{
				{
					Name: "golangci-lint",
				},
			},
		},
	}
}

func (a Action) Execute() (err error) {
	ctx, err := a.Sdk.ModuleAction(nil)
	if err != nil {
		return err
	}

	if ctx.Module.BuildSystem == "gomod" {
		_, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: `golangci-lint run ` + ctx.Config.DebugFlag("bin-golangci-lint", "-v ") + `--sort-results --issues-exit-code 1`,
			WorkDir: ctx.Module.ModuleDir,
		})
		if err != nil {
			return err
		}
	} else {
		return errors.New("build system " + ctx.Module.BuildSystem + " is not supported")
	}

	return nil
}
