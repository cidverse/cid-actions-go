package techdocs

import (
	"testing"

	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/cidverse/cid-sdk-go/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTechdocsBuild(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On(`ModuleAction`, mock.Anything).Return(TechdocsTestData(false), nil)
	sdk.On(`ExecuteCommand`, cidsdk.ExecuteCommandRequest{
		Command: `techdocs-cli generate --source-dir /my-project/docs --output-dir /my-project/.tmp/html --no-docker --etag ${NCI_COMMIT_SHA}`,
		WorkDir: `/my-project`,
	}).Return(&cidsdk.ExecuteCommandResponse{Code: 0}, nil)
	sdk.On("TARCreate", "/my-project/.tmp/html", "/my-project/.tmp/docs.tar").Return(nil)
	sdk.On("ArtifactUpload", cidsdk.ArtifactUploadRequest{
		Module:        "my-module",
		File:          "/my-project/.tmp/docs.tar",
		Type:          "html",
		Format:        "tar",
		FormatVersion: "",
		ExtractFile:   true,
	}).Return(nil)

	action := BuildAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
