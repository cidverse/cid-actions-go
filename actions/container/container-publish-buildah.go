package container

import (
	"fmt"
	"path"
	"strings"

	cidsdk "github.com/cidverse/cid-sdk-go"
)

type BuildahPublishAction struct {
	Sdk cidsdk.SDKClient
}

type BuildahPublishConfig struct {
	AlwaysPublishManifest bool
}

func (a BuildahPublishAction) Execute() error {
	cfg := BuildahPublishConfig{AlwaysPublishManifest: true}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
	}

	// target image reference
	ociDir := path.Join(ctx.Config.ArtifactDir, ctx.Module.Slug, "oci-image")
	imageRefFile := path.Join(ociDir, "image.txt")
	imageRef, err := a.Sdk.FileRead(imageRefFile)
	if err != nil {
		return fmt.Errorf("failed to parse image reference from %s", err.Error())
	}

	// for each container archive
	files, err := a.Sdk.FileList(cidsdk.FileRequest{Directory: path.Join(ctx.Config.ArtifactDir, ctx.Module.Slug, "oci-image"), Extensions: []string{".tar"}})
	if err != nil {
		return fmt.Errorf("failed to list oci archive files: %s", err.Error())
	}
	if len(files) == 0 {
		_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "warn", Message: "can't publish image, no oci archives found", Context: map[string]interface{}{"repository": imageRef}})
		return nil
	}

	_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "info", Message: "publish container image", Context: map[string]interface{}{"repository": imageRef, "manifest_size": len(files)}})

	// allow to publish single images as non-manifests
	if !cfg.AlwaysPublishManifest && len(files) == 1 {
		_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: fmt.Sprintf("buildah push --format oci oci-archive:%s docker://%s", files[0].Path, imageRef),
			WorkDir: ctx.ProjectDir,
		})
		if err != nil {
			return err
		}

		return nil
	}

	// create manifest
	manifestName := strings.Replace(a.Sdk.UUID(), "-", "", -1)
	_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: fmt.Sprintf("buildah manifest create %s", manifestName),
		WorkDir: ctx.ProjectDir,
	})
	if err != nil {
		return err
	}

	// add images to manifest
	for _, file := range files {
		_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: fmt.Sprintf("buildah manifest add %s oci-archive:%s", manifestName, file.Path),
			WorkDir: ctx.ProjectDir,
		})
		if err != nil {
			return err
		}
	}

	// print manifest
	manifestFile := path.Join(ctx.Config.TempDir, "manifest.json")
	manifestContent, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command:       fmt.Sprintf("buildah manifest inspect %s", manifestName),
		WorkDir:       ctx.ProjectDir,
		CaptureOutput: true,
	})
	if err != nil {
		return err
	}
	err = a.Sdk.FileWrite(manifestFile, []byte(manifestContent.Stdout))
	if err != nil {
		return err
	}

	// store manifest
	err = a.Sdk.ArtifactUpload(cidsdk.ArtifactUploadRequest{
		Module:        ctx.Module.Slug,
		File:          manifestFile,
		Type:          "oci-image",
		Format:        "manifest",
		FormatVersion: "v2s2",
	})
	if err != nil {
		return err
	}

	// publish manifest to registry
	digestFile := path.Join(ctx.Config.TempDir, "digest.txt")
	_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: fmt.Sprintf("buildah manifest push --all --format oci --digestfile %s %s docker://%s", digestFile, manifestName, imageRef),
		WorkDir: ctx.ProjectDir,
	})
	if err != nil {
		return err
	}
	err = a.Sdk.ArtifactUpload(cidsdk.ArtifactUploadRequest{
		Module: ctx.Module.Slug,
		File:   digestFile,
		Type:   "oci-image",
		Format: "digest",
	})
	if err != nil {
		return err
	}

	return nil
}
