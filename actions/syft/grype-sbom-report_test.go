package syft

import (
	"testing"

	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/cidverse/cid-sdk-go/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGrypeSBOMReport(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On(`ModuleAction`, mock.Anything).Return(ContainerTestData(false), nil)
	sdk.On(`FileList`, mock.Anything).Return([]cidsdk.File{cidsdk.NewFile("/my-project/.dist/my-module/sbom/test.syft.json")}, nil)
	sdk.On(`ExecuteCommand`, cidsdk.ExecuteCommandRequest{
		Command: `grype --add-cpes-if-none --file /my-project/.dist/my-module/sbom-report/test.grype.json sbom:/my-project/.dist/my-module/sbom/test.syft.json`,
		WorkDir: `/my-project`,
		Env: map[string]string{
			"GRYPE_CHECK_FOR_APP_UPDATE": "false",
			"GRYPE_OUTPUT":               "json",
		},
	}).Return(nil, nil)

	action := ReportAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
