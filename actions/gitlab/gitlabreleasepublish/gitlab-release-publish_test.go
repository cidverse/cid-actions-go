package gitlabreleasepublish

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/cidverse/cid-actions-go/actions/api"
	"github.com/cidverse/cid-actions-go/pkg/core/test"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/google/go-github/v72/github"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGitLabReleasePublishWithChangelog(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ProjectActionDataV1").Return(github.Ptr(api.GetProjectGitLabActionData(false)), nil)
	sdk.On("ArtifactDownload", cidsdk.ArtifactDownloadRequest{
		ID:         "root|changelog|gitlab.changelog",
		TargetFile: "/my-project/.tmp/gitlab.changelog",
	}).Return(nil)
	sdk.On("FileRead", "/my-project/.tmp/gitlab.changelog").Return(`changes ...`, nil)
	sdk.On("ArtifactList", cidsdk.ArtifactListRequest{Query: `artifact_type == "binary"`}).Return(&[]cidsdk.ActionArtifact{}, nil)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "https://gitlab.com/api/v4/projects/123456/releases", httpmock.NewStringResponder(200, `{
   "tag_name":"v0.3",
   "description":"Super nice release",
   "name":"New release",
   "created_at":"2019-01-03T02:22:45.118Z",
   "released_at":"2019-01-03T02:22:45.118Z"
}`))

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}

func TestGitLabReleasePublishAutoChangelog(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ProjectActionDataV1").Return(github.Ptr(api.GetProjectGitLabActionData(false)), nil)
	sdk.On("ArtifactDownload", cidsdk.ArtifactDownloadRequest{
		ID:         "root|changelog|gitlab.changelog",
		TargetFile: "/my-project/.tmp/gitlab.changelog",
	}).Return(fmt.Errorf("a error of some kind"))
	sdk.On("ArtifactList", cidsdk.ArtifactListRequest{Query: `artifact_type == "binary"`}).Return(&[]cidsdk.ActionArtifact{}, nil)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "https://gitlab.com/api/v4/projects/123456/releases", httpmock.NewStringResponder(200, `{
   "tag_name":"v0.3",
   "description":"Super nice release",
   "name":"New release",
   "created_at":"2019-01-03T02:22:45.118Z",
   "released_at":"2019-01-03T02:22:45.118Z"
}`))

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}

func TestGitLabReleasePublishSelfHosted(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ProjectActionDataV1").Return(github.Ptr(api.GetProjectGitLabActionData(false)), nil)
	sdk.On("ArtifactDownload", cidsdk.ArtifactDownloadRequest{
		ID:         "root|changelog|gitlab.changelog",
		TargetFile: "/my-project/.tmp/gitlab.changelog",
	}).Return(nil)
	sdk.On("FileRead", "/my-project/.tmp/gitlab.changelog").Return(`changes ...`, nil)
	sdk.On("ArtifactList", cidsdk.ArtifactListRequest{Query: `artifact_type == "binary"`}).Return(&[]cidsdk.ActionArtifact{}, nil)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "https://gitlab.com/api/v4/projects/123456/releases", httpmock.NewStringResponder(200, `{
   "tag_name":"v0.3",
   "description":"Super nice release",
   "name":"New release",
   "created_at":"2019-01-03T02:22:45.118Z",
   "released_at":"2019-01-03T02:22:45.118Z"
}`))

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
