package syft

import (
	"fmt"
	"path"
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

	if ctx.Module.BuildSystem != string(cidsdk.BuildSystemContainer) {
		return fmt.Errorf("not supported: %s/%s", ctx.Module.BuildSystem, ctx.Module.BuildSystemSyntax)
	}

	// find container images
	files, err := a.Sdk.FileList(cidsdk.FileRequest{
		Directory:  path.Join(ctx.Config.ArtifactDir, ctx.Module.Slug, "oci-image"),
		Extensions: []string{"tar"},
	})
	if err != nil {
		return fmt.Errorf("failed to get files in directory: %s", err.Error())
	}

	// run sbom generation for each image
	for _, file := range files {
		buildEnv := make(map[string]string)

		baseName := path.Join(ctx.Config.ArtifactDir, ctx.Module.Slug, "sbom", strings.TrimSuffix(file.Name, ".tar"))
		var outputFormats []string
		outputFormats = append(outputFormats, "json="+baseName+".syft.json")
		outputFormats = append(outputFormats, "text="+baseName+".txt") // human-readable
		// outputFormats = append(outputFormats, "cyclonedx="+baseName+".cdx.xml")            // https://cyclonedx.org/specification/overview/
		// outputFormats = append(outputFormats, "cyclonedx-json="+baseName+".cdx.json")      // https://cyclonedx.org/specification/overview/
		outputFormats = append(outputFormats, "spdx-json="+baseName+".spdx.json")          // https://github.com/spdx/spdx-spec/blob/v2.2/schemas/spdx-schema.json
		outputFormats = append(outputFormats, "spdx-tag-value="+baseName+".spdx-tag.json") // https://spdx.github.io/spdx-spec/
		outputFormats = append(outputFormats, "github="+baseName+".github.json")           // A JSON report conforming to GitHub's dependency snapshot format

		buildEnv["SYFT_CHECK_FOR_APP_UPDATE"] = "false"
		buildEnv["SYFT_OUTPUT"] = strings.Join(outputFormats, ",")

		// scan
		var buildArgs []string
		buildArgs = append(buildArgs, `--scope all-layers`)
		if ctx.Config.Debug || ctx.Config.Log["bin-syft"] == "debug" {
			buildArgs = append(buildArgs, "-vv")
		}
		buildArgs = append(buildArgs, "oci-archive:"+file.Path)
		_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: `syft packages ` + strings.Join(buildArgs, " "),
			WorkDir: ctx.ProjectDir,
			Env:     buildEnv,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
