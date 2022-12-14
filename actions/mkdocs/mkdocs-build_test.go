package mkdocs

import (
	"testing"

	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/cidverse/cid-sdk-go/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMkdocsBuild(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On(`ModuleAction`, mock.Anything).Return(MKDocsTestData(false), nil)
	sdk.On(`ExecuteCommand`, cidsdk.ExecuteCommandRequest{
		Command: `pipenv sync`,
		WorkDir: `/my-project/docs`,
	}).Return(nil, nil)
	sdk.On(`ExecuteCommand`, cidsdk.ExecuteCommandRequest{
		Command: `pipenv run mkdocs build --site-dir /my-project/.dist/my-module/html`,
		WorkDir: `/my-project/docs`,
	}).Return(nil, nil)

	action := BuildAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
