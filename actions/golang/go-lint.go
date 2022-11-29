package golang

import (
	"errors"

	cidsdk "github.com/cidverse/cid-sdk-go"
)

type LintAction struct {
	Sdk cidsdk.SDKClient
}

func (a LintAction) Execute() (err error) {
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
