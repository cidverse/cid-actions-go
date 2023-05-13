package hugobuild

import (
	"fmt"
	"strings"

	cidsdk "github.com/cidverse/cid-sdk-go"
)

type Action struct {
	Sdk cidsdk.SDKClient
}

type Config struct {
}

func (a Action) Execute() (err error) {
	cfg := Config{}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
	}

	if ctx.Module.BuildSystem != string(cidsdk.BuildSystemHugo) || ctx.Module.BuildSystemSyntax != string(cidsdk.BuildSystemSyntaxDefault) {
		return fmt.Errorf("not supported: %s/%s", ctx.Module.BuildSystem, ctx.Module.BuildSystemSyntax)
	}
	outputDir := cidsdk.JoinPath(ctx.Config.ArtifactDir, ctx.Module.Slug, "html")

	var buildArgs []string
	buildArgs = append(buildArgs, "--source "+ctx.Module.ModuleDir)
	buildArgs = append(buildArgs, "--destination "+outputDir)
	if ctx.Config.Debug || ctx.Config.Log["bin-hugo"] == "debug" {
		buildArgs = append(buildArgs, "--log --verboseLog")
	}

	_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: `hugo --minify --gc ` + strings.Join(buildArgs, " "),
		WorkDir: ctx.ProjectDir,
	})
	if err != nil {
		return err
	}

	return nil
}
