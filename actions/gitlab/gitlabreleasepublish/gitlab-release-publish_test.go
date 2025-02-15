package gitlabreleasepublish

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/cidverse/cid-actions-go/actions/gitlab/gitlabcommon"
	"github.com/cidverse/cid-actions-go/pkg/core/test"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGithubReleasePublishWithChangelog(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ProjectAction", mock.Anything).Return(gitlabcommon.GitLabTestData(), nil)
	sdk.On("ArtifactDownload", cidsdk.ArtifactDownloadRequest{
		ID:         "root|changelog|github.changelog",
		TargetFile: "/my-project/.tmp/github.changelog",
	}).Return(nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: `glab release create "v1.2.0" -F /my-project/.tmp/github.changelog`,
		WorkDir: "/my-project",
		Env: map[string]string{
			"GITLAB_HOST":     "gitlab.com",
			"GITLAB_API_HOST": "gitlab.com",
			"GITLAB_TOKEN":    "",
		},
	}).Return(&cidsdk.ExecuteCommandResponse{Code: 0}, nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}

func TestGithubReleasePublishAutoChangelog(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ProjectAction", mock.Anything).Return(gitlabcommon.GitLabTestData(), nil)
	sdk.On("ArtifactDownload", cidsdk.ArtifactDownloadRequest{
		ID:         "root|changelog|github.changelog",
		TargetFile: "/my-project/.tmp/github.changelog",
	}).Return(fmt.Errorf("a error of some kind"))
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: `glab release create "v1.2.0" --notes "no release notes"`,
		WorkDir: "/my-project",
		Env: map[string]string{
			"GITLAB_HOST":     "gitlab.com",
			"GITLAB_API_HOST": "gitlab.com",
			"GITLAB_TOKEN":    "",
		},
	}).Return(&cidsdk.ExecuteCommandResponse{Code: 0}, nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}

func TestGithubReleasePublishSelfHosted(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ProjectAction", mock.Anything).Return(gitlabcommon.GitLabSelfHostedTestData(), nil)
	sdk.On("ArtifactDownload", cidsdk.ArtifactDownloadRequest{
		ID:         "root|changelog|github.changelog",
		TargetFile: "/my-project/.tmp/github.changelog",
	}).Return(nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: `glab config set skip_tls_verify true --host "gitlab.local"`,
		WorkDir: "/my-project",
	}).Return(&cidsdk.ExecuteCommandResponse{Code: 0}, nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: `glab release create "v1.2.0" -F /my-project/.tmp/github.changelog`,
		WorkDir: "/my-project",
		Env: map[string]string{
			"GITLAB_HOST":     "gitlab.local",
			"GITLAB_API_HOST": "gitlab.local",
			"GITLAB_TOKEN":    "",
		},
	}).Return(&cidsdk.ExecuteCommandResponse{Code: 0}, nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
