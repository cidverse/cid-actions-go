package helm

import (
	"fmt"

	cidsdk "github.com/cidverse/cid-sdk-go"
)

type PublishNexusAction struct {
	Sdk cidsdk.SDKClient
}

type PublishNexusConfig struct {
	NexusURL        string `json:"nexus_url" env:"HELM_NEXUS_URL"`
	NexusRepository string `json:"nexus_repository" env:"HELM_NEXUS_REPOSITORY"`
	NexusUsername   string `json:"nexus_username" env:"HELM_NEXUS_USERNAME"`
	NexusPassword   string `json:"nexus_password" env:"HELM_NEXUS_PASSWORD"`
}

func (a PublishNexusAction) Execute() (err error) {
	cfg := PublishNexusConfig{}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
	}

	// find charts
	artifacts, err := a.Sdk.ArtifactList(cidsdk.ArtifactListRequest{Query: `artifact_type == "helm-chart" && format == "tgz"`})
	if err != nil {
		return fmt.Errorf("failed to query artifacts: %s", err.Error())
	}

	// publish
	_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "info", Message: "uploading charts to nexus", Context: map[string]interface{}{"count": len(*artifacts), "nexus": cfg.NexusURL, "nexus_repo": cfg.NexusRepository}})
	for _, artifact := range *artifacts {
		_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "info", Message: "uploading chart", Context: map[string]interface{}{"chart": artifact.Name}})

		// download
		chartArchive := cidsdk.JoinPath(ctx.Config.TempDir, artifact.Name)
		err = a.Sdk.ArtifactDownload(cidsdk.ArtifactDownloadRequest{
			ID:         artifact.ID,
			TargetFile: chartArchive,
		})
		if err != nil {
			return fmt.Errorf("failed to load artifact with id %s: %s", artifact.ID, err.Error())
		}

		// upload
		_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "info", Message: "uploading chart to nexus", Context: map[string]interface{}{"chart": artifact.Name}})
		endpoint := cfg.NexusURL + "/service/rest/v1/components?repository=" + cfg.NexusRepository
		status, response := UploadChart(endpoint, cfg.NexusUsername, cfg.NexusPassword, chartArchive)
		if status < 200 || status >= 300 {
			_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "warn", Message: "failed to upload chart", Context: map[string]interface{}{"chart": artifact.Name, "status": status, "response": string(response)}})
			return fmt.Errorf("failed to publish chart %s: status: %d, response: %s", artifact.Name, status, string(response))
		}
	}

	return nil
}
