package techdocs

import (
	"fmt"
	"os"

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

	if ctx.Module.BuildSystem != string(cidsdk.BuildSystemMkdocs) || ctx.Module.BuildSystemSyntax != string(cidsdk.MkdocsTechdocs) {
		return fmt.Errorf("not supported: %s/%s", ctx.Module.BuildSystem, ctx.Module.BuildSystemSyntax)
	}
	outputDir := cidsdk.JoinPath(ctx.Config.TempDir, "html")
	outputFile := cidsdk.JoinPath(ctx.Config.TempDir, "docs.tar")
	_ = os.MkdirAll(outputDir, os.ModePerm)

	buildResult, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: `techdocs-cli generate --source-dir ` + ctx.Module.ModuleDir + ` --output-dir ` + outputDir + ` --no-docker --etag ${NCI_COMMIT_SHA}`,
		WorkDir: ctx.ProjectDir,
	})
	if err != nil {
		return err
	} else if buildResult.Code != 0 {
		return fmt.Errorf("techdocs-cli generate failed, exit code %d", buildResult.Code)
	}

	// create zip
	err = a.Sdk.TARCreate(outputDir, outputFile)
	if err != nil {
		return err
	}

	// store result
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
