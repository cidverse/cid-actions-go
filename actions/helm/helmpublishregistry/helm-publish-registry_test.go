package helmpublishregistry

import (
	"testing"

	"github.com/cidverse/cid-actions-go/actions/helm/helmcommon"
	"github.com/cidverse/cid-actions-go/pkg/core/test"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/stretchr/testify/assert"
)

func TestHelmPublishRegistry(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ModuleActionDataV1").Return(helmcommon.GetHelmTestData(map[string]string{
		"HELM_OCI_REPOSITORY": "localhost:5000/helm-charts",
	}, false), nil)
	sdk.On("ArtifactList", cidsdk.ArtifactListRequest{Query: `artifact_type == "helm-chart" && format == "tgz"`}).Return(&[]cidsdk.ActionArtifact{
		{
			ID:     "root/helm-chart/mychart.tgz",
			Module: "root",
			Type:   "helm-chart",
			Name:   "mychart.tgz",
			Format: "tgz",
		},
	}, nil)
	sdk.On("ArtifactDownload", cidsdk.ArtifactDownloadRequest{
		ID:         "root/helm-chart/mychart.tgz",
		TargetFile: ".tmp/mychart.tgz",
	}).Return(nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: `helm push .tmp/mychart.tgz oci://localhost:5000/helm-charts`,
		WorkDir: "/my-project",
	}).Return(&cidsdk.ExecuteCommandResponse{Code: 0}, nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
