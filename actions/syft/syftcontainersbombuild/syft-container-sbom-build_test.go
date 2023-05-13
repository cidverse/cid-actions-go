package syftcontainersbombuild

import (
	"testing"

	"github.com/cidverse/cid-actions-go/pkg/core/test"
	"github.com/cidverse/cid-actions-go/testdata"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSyftSBOMScan(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On(`ModuleAction`, mock.Anything).Return(testdata.ModuleDockerfile(), nil)
	sdk.On(`FileList`, mock.Anything).Return([]cidsdk.File{cidsdk.NewFile("/my-project/.dist/my-project/oci-image/linux_amd64.tar")}, nil)
	sdk.On(`ExecuteCommand`, cidsdk.ExecuteCommandRequest{
		Command: `syft packages --quiet --scope all-layers oci-archive:/my-project/.dist/my-project/oci-image/linux_amd64.tar`,
		WorkDir: `/my-project`,
		Env: map[string]string{
			"SYFT_CHECK_FOR_APP_UPDATE": "false",
			"SYFT_OUTPUT":               "json=/my-project/.tmp/my-project/linux_amd64.syft.json,text=/my-project/.tmp/my-project/linux_amd64.txt,spdx-json=/my-project/.tmp/my-project/linux_amd64.spdx.json",
		},
	}).Return(nil, nil)
	sdk.On("ArtifactUpload", cidsdk.ArtifactUploadRequest{
		Module:        "my-project",
		File:          "/my-project/.tmp/my-project/linux_amd64.syft.json",
		Type:          "report",
		Format:        "container-sbom",
		FormatVersion: "syft-json",
	}).Return(nil)
	sdk.On("ArtifactUpload", cidsdk.ArtifactUploadRequest{
		Module:        "my-project",
		File:          "/my-project/.tmp/my-project/linux_amd64.spdx.json",
		Type:          "report",
		Format:        "container-sbom",
		FormatVersion: "spdx-json",
	}).Return(nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
