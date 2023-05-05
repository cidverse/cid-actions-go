package github

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/cidverse/cid-actions-go/actions/api"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/cidverse/cid-sdk-go/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGithubReleasePublishWithChangelog(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On("ProjectAction", mock.Anything).Return(api.GetProjectActionData(false), nil)
	sdk.On("ArtifactDownload", cidsdk.ArtifactDownloadRequest{
		Type:       "changelog",
		Name:       "github.changelog",
		TargetFile: "/my-project/.tmp/github.changelog",
	}).Return(nil)
	sdk.On("ArtifactList", cidsdk.ArtifactListRequest{Query: `artifact_type == "binary"`}).Return(&[]cidsdk.ActionArtifact{
		{
			BuildID: "0",
			JobID:   "0",
			ID:      "my-module|binary|linux_amd64",
			Module:  "my-module",
			Name:    "linux_amd64",
			Type:    "binary",
		},
	}, nil)
	sdk.On("ArtifactDownload", cidsdk.ArtifactDownloadRequest{
		ID:         "my-module|binary|linux_amd64",
		TargetFile: "/my-project/.tmp/linux_amd64",
	}).Return(nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: `gh release create "v1.2.0" --verify-tag -F "/my-project/.tmp/github.changelog" '/my-project/.tmp/linux_amd64#linux_amd64'`,
		WorkDir: "/my-project",
		Env: map[string]string{
			"GH_TOKEN": "",
		},
	}).Return(&cidsdk.ExecuteCommandResponse{Code: 0}, nil)

	action := PublishAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}

func TestGithubReleasePublishAutoChangelog(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On("ProjectAction", mock.Anything).Return(api.GetProjectActionData(false), nil)
	sdk.On("ArtifactDownload", cidsdk.ArtifactDownloadRequest{
		Type:       "changelog",
		Name:       "github.changelog",
		TargetFile: "/my-project/.tmp/github.changelog",
	}).Return(fmt.Errorf("a error of some kind"))
	sdk.On("ArtifactList", cidsdk.ArtifactListRequest{Query: `artifact_type == "binary"`}).Return(&[]cidsdk.ActionArtifact{
		{
			BuildID: "0",
			JobID:   "0",
			ID:      "my-module|binary|linux_amd64",
			Module:  "my-module",
			Name:    "linux_amd64",
			Type:    "binary",
		},
	}, nil)
	sdk.On("ArtifactDownload", cidsdk.ArtifactDownloadRequest{
		ID:         "my-module|binary|linux_amd64",
		TargetFile: "/my-project/.tmp/linux_amd64",
	}).Return(nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: `gh release create "v1.2.0" --verify-tag --generate-notes '/my-project/.tmp/linux_amd64#linux_amd64'`,
		WorkDir: "/my-project",
		Env: map[string]string{
			"GH_TOKEN": "",
		},
	}).Return(&cidsdk.ExecuteCommandResponse{Code: 0}, nil)

	action := PublishAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
