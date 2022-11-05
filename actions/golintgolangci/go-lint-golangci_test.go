package golintgolangci

import (
	"github.com/cidverse/cid-actions-go/actions/api"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/cidverse/cid-sdk-go/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGoMod(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On("PrepareAction", nil).Return(api.GetGoModTestData(false), nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "golangci-lint run --sort-results --issues-exit-code 1",
		WorkDir: "/my-project",
	}).Return(nil, nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}

func TestGoModDebug(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On("PrepareAction", nil).Return(api.GetGoModTestData(true), nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "golangci-lint run -v --sort-results --issues-exit-code 1",
		WorkDir: "/my-project",
	}).Return(nil, nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
