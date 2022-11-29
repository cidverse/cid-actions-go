package mkdocs

import (
	"testing"

	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/cidverse/cid-sdk-go/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTechdocsStart(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On(`ModuleAction`, mock.Anything).Return(MKDocsTestData(false), nil)
	sdk.On(`ExecuteCommand`, cidsdk.ExecuteCommandRequest{
		Command: `mkdocs serve --dev-addr localhost:7600 --watch /my-project/docs`,
		WorkDir: `/my-project/docs`,
	}).Return(nil, nil)

	action := StartAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
