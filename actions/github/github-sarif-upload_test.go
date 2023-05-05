package github

import (
	_ "embed"
	"testing"

	"github.com/cidverse/cid-actions-go/actions/api"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/cidverse/cid-sdk-go/mocks"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGithubSarifUpload(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On("ProjectAction", mock.Anything).Return(api.GetProjectActionData(false), nil)
	sdk.On("ArtifactList", cidsdk.ArtifactListRequest{Query: `artifact_type == "report" && format == "sarif" && format_version == "2.1.0"`}).Return(&[]cidsdk.ActionArtifact{
		{
			BuildID:       "0",
			JobID:         "0",
			Module:        "root",
			Name:          "report.sarif",
			Type:          "report",
			Format:        "sarif",
			FormatVersion: "2.1.0",
		},
	}, nil)
	sdk.On("ArtifactDownloadByteArray", cidsdk.ArtifactDownloadByteArrayRequest{
		Module: "root",
		Type:   "report",
		Name:   "report.sarif",
	}).Return([]byte("content"), nil)

	// http mock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "https://api.github.com/repos/cidverse/normalizeci/code-scanning/sarifs", httpmock.NewStringResponder(200, `{"id": "47177e22-5596-11eb-80a1-c1e54ef945c6","url": "https://api.github.com/repos/octocat/hello-world/code-scanning/sarifs/47177e22-5596-11eb-80a1-c1e54ef945c6"}`))

	// run action
	action := SarifUploadAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
