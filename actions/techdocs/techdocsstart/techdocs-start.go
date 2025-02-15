package techdocsstart

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
		Name:        "techdocs-start",
		Description: "Starts the techdocs development server locally.",
		Category:    "dev",
		Scope:       cidsdk.ActionScopeModule,
		Rules: []cidsdk.ActionRule{
			{
				Type:       "cel",
				Expression: `MODULE_BUILD_SYSTEM == "mkdocs" && MODULE_BUILD_SYSTEM_SYNTAX == "mkdocs-techdocs"`,
			},
		},
		Access: cidsdk.ActionAccess{
			Environment: []cidsdk.ActionAccessEnv{},
			Executables: []cidsdk.ActionAccessExecutable{
				{
					Name: "techdocs-cli",
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

	if ctx.Module.BuildSystem != string(cidsdk.BuildSystemMkdocs) || ctx.Module.BuildSystemSyntax != string(cidsdk.BuildSystemSyntaxMkdocsTechdocs) {
		return fmt.Errorf("not supported: %s/%s", ctx.Module.BuildSystem, ctx.Module.BuildSystemSyntax)
	}

	var startArgs []string
	startArgs = append(startArgs, "--no-docker")
	startArgs = append(startArgs, "--mkdocs-port "+strconv.Itoa(cfg.Port))
	if ctx.Config.Debug || ctx.Config.Log["bin-techdocs-cli"] == "debug" {
		startArgs = append(startArgs, "-v")
	}

	_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: fmt.Sprintf(`techdocs-cli serve %s`, strings.Join(startArgs, " ")),
		WorkDir: ctx.Module.ModuleDir,
		Ports:   []int{cfg.Port},
	})
	if err != nil {
		return err
	}

	return nil
}
