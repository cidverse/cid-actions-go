package helm

import (
	"fmt"

	cidsdk "github.com/cidverse/cid-sdk-go"
)

type PublishRegistryAction struct {
	Sdk cidsdk.SDKClient
}

type PublishRegistryConfig struct {
	OCIRepository string `json:"helm_oci_repository" env:"HELM_OCI_REPOSITORY"`
}

func (a PublishRegistryAction) Execute() (err error) {
	cfg := PublishRegistryConfig{}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
	}

	// find charts
	artifacts, err := a.Sdk.ArtifactList(cidsdk.ArtifactListRequest{Query: `artifact_type == "helm-chart" && format == "tgz"`})
	if err != nil {
		return err
	}

	// publish
	for _, artifact := range *artifacts {
		_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "info", Message: "uploading chart", Context: map[string]interface{}{"chart": artifact.Name}})

		// get chart archive
		chartArchive := cidsdk.JoinPath(ctx.Config.TempDir, artifact.Name)
		err = a.Sdk.ArtifactDownload(cidsdk.ArtifactDownloadRequest{
			ID:         artifact.ID,
			TargetFile: chartArchive,
		})
		if err != nil {
			return fmt.Errorf("failed to load artifact with id " + artifact.ID)
		}

		_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "info", Message: "uploading chart to registry", Context: map[string]interface{}{"chart": artifact.Name}})
		_, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: fmt.Sprintf(`helm push %s oci://%s`, chartArchive, cfg.OCIRepository),
			WorkDir: ctx.Module.ProjectDir,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
