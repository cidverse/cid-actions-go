package grypesbomreport

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

func (a Action) Metadata() cidsdk.ActionMetadata {
	return cidsdk.ActionMetadata{
		Name:        "grype-container-sbom-report",
		Description: "Generate a SBOM report for sbom filese",
		Category:    "security",
		Scope:       cidsdk.ActionScopeModule,
		Rules: []cidsdk.ActionRule{
			{
				Type:       "cel",
				Expression: `MODULE_BUILD_SYSTEM == "container"`,
			},
		},
		Access: cidsdk.ActionAccess{
			Environment: []cidsdk.ActionAccessEnv{},
			Executables: []cidsdk.ActionAccessExecutable{
				{
					Name: "grype",
				},
			},
		},
	}
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
		Directory:  cidsdk.JoinPath(ctx.Config.ArtifactDir, ctx.Module.Slug, "sbom"),
		Extensions: []string{".syft.json"},
	})
	if err != nil {
		return fmt.Errorf("failed to get files in directory: %s", err.Error())
	}

	// run sbom generation for each image
	for _, file := range files {
		outputFile := cidsdk.JoinPath(ctx.Config.ArtifactDir, ctx.Module.Slug, "sbom-report", file.NameShort+".grype.json")

		buildEnv := make(map[string]string)
		buildEnv["GRYPE_CHECK_FOR_APP_UPDATE"] = "false"
		buildEnv["GRYPE_OUTPUT"] = "json"
		// multiple formats blocked by https://github.com/anchore/grype/issues/648
		// var outputFormats []string
		// outputFormats = append(outputFormats, "json="+file.NameShort+".grype.json")
		// outputFormats = append(outputFormats, "table="+file.NameShort+".grype.txt")   // human-readable
		// outputFormats = append(outputFormats, "cyclonedx="+file.NameShort+".cdx.xml") // https://cyclonedx.org/specification/overview/
		// outputFormats = append(outputFormats, "sarif="+file.NameShort+".sarif")       // https://docs.oasis-open.org/sarif/sarif/v2.1.0/csprd01/sarif-v2.1.0-csprd01.html
		// buildEnv["GRYPE_OUTPUT"] = strings.Join(outputFormats, ",")

		// scan
		var buildArgs []string
		buildArgs = append(buildArgs, `--file `+outputFile)
		buildArgs = append(buildArgs, `sbom:`+file.Path)
		_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: `grype --add-cpes-if-none ` + strings.Join(buildArgs, " "),
			WorkDir: ctx.ProjectDir,
			Env:     buildEnv,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
