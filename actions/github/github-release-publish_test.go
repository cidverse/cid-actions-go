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
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: `gh release create "v1.2.0" -F "/my-project/.tmp/github.changelog"`,
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
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: `gh release create "v1.2.0" --generate-notes`,
		WorkDir: "/my-project",
		Env: map[string]string{
			"GH_TOKEN": "",
		},
	}).Return(&cidsdk.ExecuteCommandResponse{Code: 0}, nil)

	action := PublishAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
