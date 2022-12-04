package helm

import (
	"testing"

	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/cidverse/cid-sdk-go/mocks"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHelmPublishNexus(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On("ModuleAction", &PublishNexusConfig{}).Return(GetHelmTestData(false), nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*PublishNexusConfig)
		arg.NexusURL = "https://localhost:9999"
		arg.NexusRepository = "dummy"
		arg.NexusUsername = "admin"
		arg.NexusPassword = "admin"
	})
	sdk.On("FileList", cidsdk.FileRequest{Directory: ".dist/my-package/helm-charts", Extensions: []string{".tgz"}}).Return([]cidsdk.File{cidsdk.NewFile(".dist/my-package/helm-charts/my-chart.tgz")}, nil)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "https://localhost:9999/service/rest/v1/components?repository=dummy", httpmock.NewStringResponder(200, ``))

	action := PublishNexusAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
