package mkdocs

import (
	"fmt"
	"os"
	"strings"

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

	if ctx.Module.BuildSystem != string(cidsdk.BuildSystemMkdocs) || ctx.Module.BuildSystemSyntax != string(cidsdk.BuildSystemSyntaxDefault) {
		return fmt.Errorf("not supported: %s/%s", ctx.Module.BuildSystem, ctx.Module.BuildSystemSyntax)
	}
	outputDir := cidsdk.JoinPath(ctx.Config.TempDir, "html")
	outputFile := cidsdk.JoinPath(ctx.Config.TempDir, "docs.tar")
	_ = os.MkdirAll(outputDir, os.ModePerm)

	// install
	_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: `pipenv sync`,
		WorkDir: ctx.Module.ModuleDir,
	})
	if err != nil {
		return err
	}

	// mkdocs
	var mkdocsArgs []string
	mkdocsArgs = append(mkdocsArgs, "--site-dir "+outputDir)
	_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: `pipenv run mkdocs build ` + strings.Join(mkdocsArgs, " "),
		WorkDir: ctx.Module.ModuleDir,
	})
	if err != nil {
		return err
	}

	// create zip
	err = a.Sdk.TARCreate(outputDir, outputFile)
	if err != nil {
		return err
	}

	// store result
	_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "info", Message: "Uploading artifact...", Context: map[string]interface{}{"file": outputFile}})
	err = a.Sdk.ArtifactUpload(cidsdk.ArtifactUploadRequest{
		File:        outputFile,
		Module:      ctx.Module.Slug,
		Type:        "html",
		Format:      "tar",
		ExtractFile: true,
	})
	if err != nil {
		return err
	}

	return nil
}
