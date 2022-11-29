package mkdocs

import (
	"fmt"
	"strconv"
	"strings"

	cidsdk "github.com/cidverse/cid-sdk-go"
)

type StartAction struct {
	Sdk cidsdk.SDKClient
}

type StartConfig struct {
	Port int
}

func (a StartAction) Execute() (err error) {
	cfg := StartConfig{Port: 7600}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
	}

	if ctx.Module.BuildSystem != string(cidsdk.BuildSystemMkdocs) || ctx.Module.BuildSystemSyntax != string(cidsdk.BuildSystemSyntaxDefault) {
		return fmt.Errorf("not supported: %s/%s", ctx.Module.BuildSystem, ctx.Module.BuildSystemSyntax)
	}

	var startArgs []string
	startArgs = append(startArgs, "--dev-addr localhost:"+strconv.Itoa(cfg.Port))
	startArgs = append(startArgs, "--watch "+ctx.Module.ModuleDir)

	if ctx.Config.Debug || ctx.Config.Log["bin-mkdocs-cli"] == "debug" {
		startArgs = append(startArgs, "-v")
	}

	_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: `mkdocs serve ` + strings.Join(startArgs, " "),
		WorkDir: ctx.Module.ModuleDir,
	})
	if err != nil {
		return err
	}

	return nil
}
