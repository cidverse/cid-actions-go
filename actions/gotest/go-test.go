package gotest

import (
	"errors"
	cidsdk "github.com/cidverse/cid-sdk-go"
)

type Action struct {
	Sdk cidsdk.SDKClient
}

func (a Action) Execute() (err error) {
	ctx, err := a.Sdk.PrepareAction(nil)
	if err != nil {
		return err
	}

	if ctx.Module.BuildSystem == "gomod" {
		_, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: "go test -vet off " + ctx.Config.DebugFlag("bin-go", "-v ") + "-cover -covermode=count ./...",
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
