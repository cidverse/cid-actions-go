package gotest

import (
	"errors"
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
		Command: "go test -vet off -cover -covermode=count ./...",
		WorkDir: "/my-project",
	}).Return(nil, nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}

func TestDebug(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On("PrepareAction", nil).Return(api.GetGoModTestData(true), nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "go test -vet off -v -cover -covermode=count ./...",
		WorkDir: "/my-project",
	}).Return(nil, nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}

func TestUnsupported(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On("PrepareAction", nil).Return(api.GetUnknownTestData(false), nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.Error(t, errors.New("build system unknown is not supported"), err)
}
