package helmbuild

import (
	"fmt"
	"path/filepath"

	"github.com/cidverse/cid-actions-go/actions/helm/helmcommon"
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

	// globals
	chartsDir := cidsdk.JoinPath(ctx.Config.TempDir, "helm-charts")
	chartName := filepath.Base(ctx.Module.ModuleDir)

	if ctx.Module.BuildSystem == string(cidsdk.BuildSystemHelm) {
		// restore the charts/ directory based on the Chart.lock file
		_, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: `helm dependency build .`,
			WorkDir: ctx.Module.ModuleDir,
		})
		if err != nil {
			return err
		}

		// parse chart
		chartFile := cidsdk.JoinPath(ctx.Module.ModuleDir, "Chart.yaml")
		chartFileContent, err := a.Sdk.FileRead(chartFile)
		if err != nil {
			return fmt.Errorf("failed to read chart file: %s", err.Error())
		}
		chart, err := helmcommon.ParseChart([]byte(chartFileContent))
		if err != nil {
			return err
		}
		_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "info", Message: "building chart", Context: map[string]interface{}{"chart": chart.Name}})

		// version
		chartVersion := chart.Version
		if chartVersion == "0.0.0" {
			chartVersion = ctx.Env["NCI_COMMIT_REF_NAME"]
		}

		// package
		_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: `helm package . --version ` + chartVersion + ` --destination ` + chartsDir,
			WorkDir: ctx.Module.ModuleDir,
		})
		if err != nil {
			return err
		}

		// update index
		_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: `helm repo index ` + chartsDir,
			WorkDir: ctx.Module.ModuleDir,
		})
		if err != nil {
			return err
		}

		// upload charts
		err = a.Sdk.ArtifactUpload(cidsdk.ArtifactUploadRequest{
			File:   cidsdk.JoinPath(chartsDir, fmt.Sprintf("%s-%s.tgz", chartName, chartVersion)),
			Module: ctx.Module.Slug,
			Type:   "helm-chart",
			Format: "tgz",
		})
		if err != nil {
			return err
		}
	}

	return nil
}
