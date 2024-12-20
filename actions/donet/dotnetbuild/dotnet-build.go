package dotnetbuild

import (
	"fmt"

	cidsdk "github.com/cidverse/cid-sdk-go"
)

type Action struct {
	Sdk cidsdk.SDKClient
}

type Config struct {
}

func (a Action) Execute() error {
	cfg := Config{}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
	}

	// restore
	restoreResult, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: fmt.Sprintf(`dotnet restore --disable-parallel`),
		WorkDir: ctx.Module.ModuleDir,
	})
	if err != nil {
		return err
	} else if restoreResult.Code != 0 {
		return fmt.Errorf("dotnet restore failed, exit code %d", restoreResult.Code)
	}

	// build
	buildResult, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: fmt.Sprintf(`dotnet publish --configuration Release`),
		WorkDir: ctx.Module.ModuleDir,
	})
	if err != nil {
		return err
	} else if buildResult.Code != 0 {
		return fmt.Errorf("dotnet build failed, exit code %d", buildResult.Code)
	}

	// store result
	/*
		err = a.Sdk.ArtifactUpload(cidsdk.ArtifactUploadRequest{
			File:   outputFile,
			Module: ctx.Module.Slug,
			Type:   "binary",
			Format: "go",
		})
		if err != nil {
			return err
		}
	*/

	return nil
}
