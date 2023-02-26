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

	// globals
	chartArtifactDir := cidsdk.JoinPath(ctx.Config.ArtifactDir, ctx.Module.Slug, "helm-charts")

	// publish
	files, err := a.Sdk.FileList(cidsdk.FileRequest{Directory: chartArtifactDir, Extensions: []string{".tgz"}})
	if err != nil {
		return fmt.Errorf("failed to find any charts in artifact directory: %s", err.Error())
	}
	_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "info", Message: "uploading charts to nexus", Context: map[string]interface{}{"count": len(files), "nexus": cfg.NexusURL, "nexus_repo": cfg.NexusRepository}})

	for _, file := range files {
		endpoint := cfg.NexusURL + "/service/rest/v1/components?repository=" + cfg.NexusRepository
		_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "info", Message: "uploading chart", Context: map[string]interface{}{"chart": file.Name}})

		status, response := UploadChart(endpoint, cfg.NexusUsername, cfg.NexusPassword, file.Path)
		if status < 200 || status >= 300 {
			_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "warn", Message: "failed to upload chart", Context: map[string]interface{}{"chart": file.Name, "status": status, "response": string(response)}})
			return fmt.Errorf("failed to publish chart %s: status: %d, response: %s", file.Name, status, string(response))
		}
	}

	return nil
}
