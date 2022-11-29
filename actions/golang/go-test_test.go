package golang

import (
	"errors"
	"testing"

	"github.com/cidverse/cid-actions-go/actions/api"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/cidverse/cid-sdk-go/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGoModTest(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On("ModuleAction", nil).Return(api.GetGoModTestData(false), nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "go test -vet off -cover -covermode=count ./...",
		WorkDir: "/my-project",
	}).Return(nil, nil)

	action := TestAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}

func TestDebugTest(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On("ModuleAction", nil).Return(api.GetGoModTestData(true), nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "go test -vet off -v -cover -covermode=count ./...",
		WorkDir: "/my-project",
	}).Return(nil, nil)

	action := TestAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}

func TestUnsupportedTest(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On("ModuleAction", nil).Return(api.GetUnknownTestData(false), nil)

	action := TestAction{Sdk: sdk}
	err := action.Execute()
	assert.Error(t, errors.New("build system unknown is not supported"), err)
}
