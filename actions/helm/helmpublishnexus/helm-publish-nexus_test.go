package helmpublishnexus

import (
	"testing"

	"github.com/cidverse/cid-actions-go/actions/helm/helmcommon"
	"github.com/cidverse/cid-actions-go/pkg/core/test"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHelmPublishNexus(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ModuleAction", &PublishNexusConfig{}).Return(helmcommon.GetHelmTestData(false), nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*PublishNexusConfig)
		arg.NexusURL = "https://localhost:9999"
		arg.NexusRepository = "dummy"
		arg.NexusUsername = "admin"
		arg.NexusPassword = "admin"
	})
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

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "https://localhost:9999/service/rest/v1/components?repository=dummy", httpmock.NewStringResponder(200, ``))

	action := PublishNexusAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
