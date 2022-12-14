package container

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	cidsdk "github.com/cidverse/cid-sdk-go"
)

type BuildahBuildAction struct {
	Sdk cidsdk.SDKClient
}

type BuildahBuildConfig struct {
	NoCache bool `json:"no-cache"`
	Squash  bool `json:"squash"`
	Rebuild bool `json:"rebuild"`
}

func (a BuildahBuildAction) Execute() error {
	cfg := BuildahBuildConfig{NoCache: false, Squash: true, Rebuild: false}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
	}

	imageReference := getFullImage(ctx.Env["NCI_CONTAINERREGISTRY_HOST"], ctx.Env["NCI_CONTAINERREGISTRY_REPOSITORY"], ctx.Env["NCI_CONTAINERREGISTRY_TAG"])
	for _, discovery := range ctx.Module.Discovery {
		containerFile := discovery.File
		containerFileContent, _ := a.Sdk.FileRead(containerFile)
		_ = os.MkdirAll(path.Join(ctx.Config.ArtifactDir, ctx.Module.Slug, "oci-image"), os.ModePerm)

		if ctx.Module.BuildSystem == string(cidsdk.BuildSystemContainer) && ctx.Module.BuildSystemSyntax == string(cidsdk.ContainerFile) {
			platforms := getDockerfileTargetPlatforms(containerFileContent)
			imageReference = getDockerfileTargetImageWithVersion(containerFileContent, imageReference)

			// build each image and add to manifest
			for _, platform := range platforms {
				containerArchiveFile := path.Join(ctx.Config.ArtifactDir, ctx.Module.Slug, "oci-image", platform.Platform("_")+".tar")
				_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "info", Message: "build container image", Context: map[string]interface{}{"module": ctx.Module.Name, "platform": platform.Platform("/"), "tag": imageReference}})

				var buildArgs []string
				buildArgs = append(buildArgs, "--platform "+platform.Platform("/"))
				buildArgs = append(buildArgs, "-f "+filepath.Base(containerFile))
				buildArgs = append(buildArgs, fmt.Sprintf("-t oci-archive:%s", strings.TrimPrefix(containerArchiveFile, ctx.ProjectDir+"/"))) // requires a relative path for some reason
				buildArgs = append(buildArgs, "--layers")                                                                                     // enable layer cache

				// options
				if cfg.NoCache {
					buildArgs = append(buildArgs, "--no-cache")
				}
				if cfg.Squash {
					buildArgs = append(buildArgs, "--squash") // squash, excluding the base layer
				}

				// labels (oci annotations: https://github.com/opencontainers/image-spec/blob/main/annotations.md)
				buildArgs = append(buildArgs, `--annotation "org.opencontainers.image.source=${NCI_REPOSITORY_REMOTE}"`)
				buildArgs = append(buildArgs, `--annotation "org.opencontainers.image.created=${TIMESTAMP_RFC3339}"`)
				buildArgs = append(buildArgs, `--annotation "org.opencontainers.image.authors="`)
				buildArgs = append(buildArgs, `--annotation "org.opencontainers.image.title=`+ctx.Module.Name+`"`)
				buildArgs = append(buildArgs, `--annotation "org.opencontainers.image.description="`)

				// dynamic build-args
				if strings.Contains(containerFileContent, "ARG TARGETPLATFORM") {
					buildArgs = append(buildArgs, `--build-arg TARGETPLATFORM=`+platform.Platform("/"))
				}
				if strings.Contains(containerFileContent, "ARG TARGETOS") {
					buildArgs = append(buildArgs, `--build-arg TARGETOS=`+platform.OS)
				}
				if strings.Contains(containerFileContent, "ARG TARGETARCH") {
					buildArgs = append(buildArgs, `--build-arg TARGETARCH=`+platform.Arch)
				}
				if strings.Contains(containerFileContent, "ARG TARGETVARIANT") {
					buildArgs = append(buildArgs, `--build-arg TARGETVARIANT=`+platform.Variant)
				}

				_, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
					Command: fmt.Sprintf("buildah build %s %s", strings.Join(buildArgs, " "), ctx.Module.ModuleDir),
					WorkDir: ctx.ProjectDir,
				})
				if err != nil {
					return err
				}
			}
		} else {
			return fmt.Errorf("not supported: %s/%s", ctx.Module.BuildSystem, ctx.Module.BuildSystemSyntax)
		}
	}

	// store image ref
	err = a.Sdk.FileWrite(path.Join(ctx.Config.ArtifactDir, ctx.Module.Slug, "oci-image", "image.txt"), []byte(imageReference))
	if err != nil {
		return errors.New("failed to store file with the image ref: " + err.Error())
	}

	return nil
}
