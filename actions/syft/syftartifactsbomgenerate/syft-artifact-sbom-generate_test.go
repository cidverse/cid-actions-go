package syftartifactsbomgenerate

import (
	"testing"

	"github.com/cidverse/cid-actions-go/pkg/core/test"
	"github.com/cidverse/cid-actions-go/testdata"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSyftArtifactSBOMGenerate(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On(`ModuleAction`, mock.Anything).Return(testdata.ModuleDockerfile(), nil)
	sdk.On("ArtifactList", cidsdk.ArtifactListRequest{Query: `module == "my-project" && artifact_type == "binary"`}).Return(&[]cidsdk.ActionArtifact{
		{
			BuildID: "0",
			JobID:   "0",
			ID:      "my-project|binary|linux_amd64",
			Module:  "my-project",
			Type:    "binary",
			Name:    "linux_amd64",
		},
	}, nil)
	sdk.On("ArtifactDownload", cidsdk.ArtifactDownloadRequest{
		ID:         "my-project|binary|linux_amd64",
		TargetFile: "/my-project/.tmp/linux_amd64",
	}).Return(nil)
	sdk.On(`ExecuteCommand`, cidsdk.ExecuteCommandRequest{
		Command: `syft packages --quiet file:/my-project/.tmp/linux_amd64`,
		WorkDir: `/my-project`,
		Env: map[string]string{
			"SYFT_CHECK_FOR_APP_UPDATE": "false",
			"SYFT_OUTPUT":               "json=/my-project/.tmp/linux_amd64.syft.json,spdx-json=/my-project/.tmp/linux_amd64.spdx.json",
		},
	}).Return(nil, nil)
	sdk.On("ArtifactUpload", cidsdk.ArtifactUploadRequest{
		Module:        "my-project",
		File:          "/my-project/.tmp/linux_amd64.syft.json",
		Type:          "report",
		Format:        "artifact-sbom",
		FormatVersion: "syft-json",
	}).Return(nil)
	sdk.On("ArtifactUpload", cidsdk.ArtifactUploadRequest{
		Module:        "my-project",
		File:          "/my-project/.tmp/linux_amd64.spdx.json",
		Type:          "report",
		Format:        "artifact-sbom",
		FormatVersion: "spdx-json",
	}).Return(nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
