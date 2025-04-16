package hugostart

import (
	"fmt"
	"strconv"
	"strings"

	cidsdk "github.com/cidverse/cid-sdk-go"
)

type Action struct {
	Sdk cidsdk.SDKClient
}

type Config struct {
	Port int
}

func (a Action) Metadata() cidsdk.ActionMetadata {
	return cidsdk.ActionMetadata{
		Name:        "hugo-start",
		Description: "Starts the hugo server for local development.",
		Category:    "dev",
		Scope:       cidsdk.ActionScopeModule,
		Rules: []cidsdk.ActionRule{
			{
				Type:       "cel",
				Expression: `MODULE_BUILD_SYSTEM == "hugo"`,
			},
		},
		Access: cidsdk.ActionAccess{
			Environment: []cidsdk.ActionAccessEnv{},
			Executables: []cidsdk.ActionAccessExecutable{
				{
					Name: "hugo",
				},
			},
		},
	}
}

func (a Action) Execute() (err error) {
	cfg := Config{Port: 7600}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
	}

	if ctx.Module.BuildSystem != string(cidsdk.BuildSystemHugo) || ctx.Module.BuildSystemSyntax != string(cidsdk.BuildSystemSyntaxDefault) {
		return fmt.Errorf("not supported: %s/%s", ctx.Module.BuildSystem, ctx.Module.BuildSystemSyntax)
	}

	var startArgs []string
	startArgs = append(startArgs, `--source `+ctx.Module.ModuleDir)
	startArgs = append(startArgs, `--minify --gc`)
	startArgs = append(startArgs, `--baseUrl /`)
	startArgs = append(startArgs, `--watch`)
	startArgs = append(startArgs, `--port `+strconv.Itoa(cfg.Port))
	if ctx.Config.Debug || ctx.Config.Log["bin-hugo"] == "debug" {
		startArgs = append(startArgs, "--debug")
	}

	_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: `hugo server ` + strings.Join(startArgs, " "),
		WorkDir: ctx.Module.ModuleDir,
		Ports:   []int{cfg.Port},
	})
	if err != nil {
		return err
	}

	return nil
}
