package syftcontainersbombuild

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

	if ctx.Module.BuildSystem != string(cidsdk.BuildSystemContainer) {
		return fmt.Errorf("not supported: %s/%s", ctx.Module.BuildSystem, ctx.Module.BuildSystemSyntax)
	}

	// find container images
	files, err := a.Sdk.FileList(cidsdk.FileRequest{
		Directory:  cidsdk.JoinPath(ctx.Config.ArtifactDir, ctx.Module.Slug, "oci-image"),
		Extensions: []string{"tar"},
	})
	if err != nil {
		return fmt.Errorf("failed to get files in directory: %s", err.Error())
	}

	// run sbom generation for each image
	for _, file := range files {
		baseName := cidsdk.JoinPath(ctx.Config.TempDir, ctx.Module.Slug, strings.TrimSuffix(file.Name, ".tar"))
		outputFormats := []string{
			"json=" + baseName + ".syft.json",      // syft-json
			"text=" + baseName + ".txt",            // human-readable
			"spdx-json=" + baseName + ".spdx.json", // https://github.com/spdx/spdx-spec/blob/v2.2/schemas/spdx-schema.json
		}

		// scan
		var buildArgs []string
		buildArgs = append(buildArgs, `--scope all-layers`)
		if ctx.Config.Debug || ctx.Config.Log["bin-syft"] == "debug" {
			buildArgs = append(buildArgs, "-vv")
		}
		buildArgs = append(buildArgs, "oci-archive:"+file.Path)
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
			File:          baseName + ".syft.json",
			Type:          "report",
			Format:        "container-sbom",
			FormatVersion: "syft-json",
		})
		if err != nil {
			return err
		}
		err = a.Sdk.ArtifactUpload(cidsdk.ArtifactUploadRequest{
			Module:        ctx.Module.Slug,
			File:          baseName + ".spdx.json",
			Type:          "report",
			Format:        "container-sbom",
			FormatVersion: "spdx-json",
		})
		if err != nil {
			return err
		}
	}

	return nil
}
