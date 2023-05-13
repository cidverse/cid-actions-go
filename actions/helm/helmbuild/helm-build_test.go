package helmbuild

import (
	"testing"

	"github.com/cidverse/cid-actions-go/actions/helm/helmcommon"
	"github.com/cidverse/cid-actions-go/pkg/core/test"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/stretchr/testify/assert"
)

func TestHelmBuild(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ModuleAction", &BuildConfig{}).Return(helmcommon.GetHelmTestData(false), nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "helm dependency build .",
		WorkDir: "/my-project/charts/mychart",
	}).Return(nil, nil)
	sdk.On("FileRead", "/my-project/charts/mychart/Chart.yaml").Return("name: mychart\nversion: 1.1.0", nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "helm package . --version 1.1.0 --destination .tmp/helm-charts",
		WorkDir: "/my-project/charts/mychart",
	}).Return(nil, nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "helm repo index .tmp/helm-charts",
		WorkDir: "/my-project/charts/mychart",
	}).Return(nil, nil)
	sdk.On("ArtifactUpload", cidsdk.ArtifactUploadRequest{
		Module: "my-package",
		File:   ".tmp/helm-charts/mychart-1.1.0.tgz",
		Type:   "helm-chart",
		Format: "tgz",
	}).Return(nil)

	action := BuildAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
