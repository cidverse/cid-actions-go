package syft

import (
	"fmt"
	"strings"

	cidsdk "github.com/cidverse/cid-sdk-go"
)

type ArtifactGenerateAction struct {
	Sdk cidsdk.SDKClient
}

type ArtifactGenerateConfig struct {
}

func (a ArtifactGenerateAction) Execute() (err error) {
	cfg := ArtifactGenerateConfig{}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
	}

	// find artifacts
	artifacts, err := a.Sdk.ArtifactList(cidsdk.ArtifactListRequest{ArtifactType: "binary", Module: ctx.Module.Slug})
	if err != nil {
		return err
	}

	// run sbom generation for each image
	for _, file := range *artifacts {
		targetFile := cidsdk.JoinPath(ctx.Config.TempDir, file.Name)
		var dlErr = a.Sdk.ArtifactDownload(cidsdk.ArtifactDownloadRequest{
			Module:     file.Module,
			Type:       string(file.Type),
			Name:       file.Name,
			TargetFile: targetFile,
		})
		if dlErr != nil {
			_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "warn", Message: "failed to retrieve sarif report", Context: map[string]interface{}{"artifact": fmt.Sprintf("%s-%s", file.Module, file.Name)}})
			continue
		}

		// format
		outputFormats := []string{
			"json=" + targetFile + ".syft.json",      // syft-json
			"spdx-json=" + targetFile + ".spdx.json", // https://github.com/spdx/spdx-spec/blob/v2.2/schemas/spdx-schema.json
		}

		// scan
		var buildArgs []string
		buildArgs = append(buildArgs, "file:"+targetFile)
		_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: `syft packages --quiet ` + strings.Join(buildArgs, " "),
			WorkDir: ctx.ProjectDir,
			Env: map[string]string{
				"SYFT_CHECK_FOR_APP_UPDATE": "false",
				"SYFT_OUTPUT":               strings.Join(outputFormats, ","),
			},
		})
		if err != nil {
			return err
		}

		// store reports
		err = a.Sdk.ArtifactUpload(cidsdk.ArtifactUploadRequest{
			Module:        ctx.Module.Slug,
			File:          targetFile + ".syft.json",
			Type:          "report",
			Format:        "artifact-sbom",
			FormatVersion: "syft-json",
		})
		if err != nil {
			return err
		}
		err = a.Sdk.ArtifactUpload(cidsdk.ArtifactUploadRequest{
			Module:        ctx.Module.Slug,
			File:          targetFile + ".spdx.json",
			Type:          "report",
			Format:        "artifact-sbom",
			FormatVersion: "spdx-json",
		})
		if err != nil {
			return err
		}
	}

	return nil
}
