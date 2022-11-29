package syft

import (
	"testing"

	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/cidverse/cid-sdk-go/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSyftSBOMBuild(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On(`ModuleAction`, mock.Anything).Return(ContainerTestData(false), nil)
	sdk.On(`FileList`, mock.Anything).Return([]cidsdk.File{cidsdk.NewFile("/my-project/.dist/my-module/oci-image/linux_amd64.tar")}, nil)
	sdk.On(`ExecuteCommand`, cidsdk.ExecuteCommandRequest{
		Command: `syft packages --scope all-layers oci-archive:/my-project/.dist/my-module/oci-image/linux_amd64.tar`,
		WorkDir: `/my-project`,
		Env: map[string]string{
			"SYFT_CHECK_FOR_APP_UPDATE": "false",
			"SYFT_OUTPUT":               "json=/my-project/.dist/my-module/sbom/linux_amd64.syft.json,text=/my-project/.dist/my-module/sbom/linux_amd64.txt,spdx-json=/my-project/.dist/my-module/sbom/linux_amd64.spdx.json,spdx-tag-value=/my-project/.dist/my-module/sbom/linux_amd64.spdx-tag.json,github=/my-project/.dist/my-module/sbom/linux_amd64.github.json",
		},
	}).Return(nil, nil)

	action := BuildAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
