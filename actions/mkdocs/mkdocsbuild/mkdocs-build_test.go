package mkdocsbuild

import (
	"testing"

	"github.com/cidverse/cid-actions-go/actions/mkdocs/mkdocscommon"
	"github.com/cidverse/cid-actions-go/pkg/core/test"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMkdocsBuild(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On(`ModuleAction`, mock.Anything).Return(mkdocscommon.MKDocsTestData(false), nil)
	sdk.On(`ExecuteCommand`, cidsdk.ExecuteCommandRequest{
		Command: `pipenv sync`,
		WorkDir: `/my-project/docs`,
	}).Return(nil, nil)
	sdk.On(`ExecuteCommand`, cidsdk.ExecuteCommandRequest{
		Command: `pipenv run mkdocs build --site-dir /my-project/.tmp/html`,
		WorkDir: `/my-project/docs`,
	}).Return(nil, nil)
	sdk.On("TARCreate", "/my-project/.tmp/html", "/my-project/.tmp/docs.tar").Return(nil)
	sdk.On("ArtifactUpload", cidsdk.ArtifactUploadRequest{
		Module:        "my-module",
		File:          "/my-project/.tmp/docs.tar",
		Type:          "html",
		Format:        "tar",
		FormatVersion: "",
		ExtractFile:   true,
	}).Return(nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
