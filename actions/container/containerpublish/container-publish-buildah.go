package containerpublish

import (
	"fmt"
	"strings"

	"github.com/cidverse/cid-actions-go/util/container"
	cidsdk "github.com/cidverse/cid-sdk-go"
)

type Action struct {
	Sdk cidsdk.SDKClient
}

type Config struct {
	AlwaysPublishManifest bool `json:"buildah_always_publish_manifest" env:"BUILDAH_ALWAYS_PUBLISH_MANIFEST"`
}

func (a Action) Execute() error {
	cfg := Config{AlwaysPublishManifest: false}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
	}

	// properties
	digestFile := cidsdk.JoinPath(ctx.Config.TempDir, "digest.txt")

	// target image reference
	ociDir := cidsdk.JoinPath(ctx.Config.ArtifactDir, ctx.Module.Slug, "oci-image")
	imageRefFile := cidsdk.JoinPath(ociDir, "image.txt")
	imageRef, err := a.Sdk.FileRead(imageRefFile)
	if err != nil {
		return fmt.Errorf("failed to parse image reference from %s", err.Error())
	}

	// for each container archive
	files, err := a.Sdk.FileList(cidsdk.FileRequest{Directory: cidsdk.JoinPath(ctx.Config.ArtifactDir, ctx.Module.Slug, "oci-image"), Extensions: []string{".tar"}})
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
		// push
		pushResult, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: fmt.Sprintf(`skopeo copy --digestfile %s oci-archive:%s docker://%s`, digestFile, files[0].Path, imageRef),
			WorkDir: ctx.ProjectDir,
		})
		if err != nil {
			return err
		} else if pushResult.Code != 0 {
			return fmt.Errorf("failed, exit code %d", pushResult.Code)
		}

		// store digest
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

	// create manifest
	manifestName := strings.Replace(a.Sdk.UUID(), "-", "", -1)
	manifestCreateResult, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: fmt.Sprintf("buildah manifest create %s", manifestName),
		WorkDir: ctx.ProjectDir,
	})
	if err != nil {
		return err
	} else if manifestCreateResult.Code != 0 {
		return fmt.Errorf("failed, exit code %d", manifestCreateResult.Code)
	}

	// add images to manifest
	for _, file := range files {
		manifestAddResult, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: fmt.Sprintf("buildah manifest add %s oci-archive:%s", manifestName, file.Path),
			WorkDir: ctx.ProjectDir,
		})
		if err != nil {
			return err
		} else if manifestAddResult.Code != 0 {
			return fmt.Errorf("failed, exit code %d", manifestAddResult.Code)
		}
	}

	// publish manifest to registry
	pushResult, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: fmt.Sprintf("buildah manifest push --all --format oci --digestfile %s %s docker://%s", digestFile, manifestName, imageRef),
		WorkDir: ctx.ProjectDir,
	})
	if err != nil {
		return err
	} else if pushResult.Code != 0 {
		return fmt.Errorf("failed, exit code %d", pushResult.Code)
	}

	// store digest
	err = a.Sdk.ArtifactUpload(cidsdk.ArtifactUploadRequest{
		Module: ctx.Module.Slug,
		File:   digestFile,
		Type:   "oci-image",
		Format: "digest",
	})
	if err != nil {
		return err
	}
	digest, err := a.Sdk.FileRead(digestFile)
	if err != nil {
		return fmt.Errorf("failed to read digest file: %s", err.Error())
	}

	// retrieve manifest
	manifestResult, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command:       fmt.Sprintf("buildah manifest inspect %s@%s", container.GetImageReferenceWithoutTag(imageRef), digest),
		WorkDir:       ctx.ProjectDir,
		CaptureOutput: true,
	})
	if err != nil {
		return err
	} else if manifestResult.Code != 0 {
		return fmt.Errorf("failed, exit code %d", manifestResult.Code)
	}

	// store manifest with correct digests
	manifestFile := cidsdk.JoinPath(ctx.Config.TempDir, "manifest.json")
	err = a.Sdk.FileWrite(manifestFile, []byte(manifestResult.Stdout))
	if err != nil {
		return err
	}

	// upload manifest
	err = a.Sdk.ArtifactUpload(cidsdk.ArtifactUploadRequest{
		Module:        ctx.Module.Slug,
		File:          manifestFile,
		Type:          "oci-image",
		Format:        "manifest",
		FormatVersion: "oci",
	})
	if err != nil {
		return err
	}

	return nil
}
