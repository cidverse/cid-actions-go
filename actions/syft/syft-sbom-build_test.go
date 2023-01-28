package syft

import (
	"testing"

	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/cidverse/cid-sdk-go/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSyftSBOMScan(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On(`ModuleAction`, mock.Anything).Return(ContainerTestData(false), nil)
	sdk.On(`FileList`, mock.Anything).Return([]cidsdk.File{cidsdk.NewFile("/my-project/.dist/my-module/oci-image/linux_amd64.tar")}, nil)
	sdk.On(`ExecuteCommand`, cidsdk.ExecuteCommandRequest{
		Command: `syft packages --quiet --scope all-layers oci-archive:/my-project/.dist/my-module/oci-image/linux_amd64.tar`,
		WorkDir: `/my-project`,
		Env: map[string]string{
			"SYFT_CHECK_FOR_APP_UPDATE": "false",
			"SYFT_OUTPUT":               "json=/my-project/.tmp/my-module/linux_amd64.syft.json,text=/my-project/.tmp/my-module/linux_amd64.txt,spdx-json=/my-project/.tmp/my-module/linux_amd64.spdx.json",
		},
	}).Return(nil, nil)
	sdk.On("ArtifactUpload", cidsdk.ArtifactUploadRequest{
		Module:        "my-module",
		File:          "/my-project/.tmp/my-module/linux_amd64.syft.json",
		Type:          "report",
		Format:        "container-sbom",
		FormatVersion: "syft-json",
	}).Return(nil)
	sdk.On("ArtifactUpload", cidsdk.ArtifactUploadRequest{
		Module:        "my-module",
		File:          "/my-project/.tmp/my-module/linux_amd64.spdx.json",
		Type:          "report",
		Format:        "container-sbom",
		FormatVersion: "spdx-json",
	}).Return(nil)

	action := ScanAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
