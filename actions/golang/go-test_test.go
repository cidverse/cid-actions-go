package golang

import (
	"errors"
	"testing"

	"github.com/cidverse/cid-actions-go/actions/api"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/cidverse/cid-sdk-go/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGoModTest(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On("ModuleAction", mock.Anything).Return(GoModTestData(false), nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "go test -vet off -cover -covermode=count -coverprofile /my-project/.dist/github-com-cidverse-my-project/go-test/cover.out ./...",
		WorkDir: "/my-project",
	}).Return(nil, nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command:       "go test -coverprofile /my-project/.dist/github-com-cidverse-my-project/go-test/cover.out -json -covermode=count ./...",
		WorkDir:       "/my-project",
		CaptureOutput: true,
	}).Return(&cidsdk.ExecuteCommandResponse{Stdout: "{}"}, nil)
	sdk.On("FileWrite", "/my-project/.dist/github-com-cidverse-my-project/go-test/cover.json", []byte("{}")).Return(nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "go tool cover -html /my-project/.dist/github-com-cidverse-my-project/go-test/cover.out -o /my-project/.dist/github-com-cidverse-my-project/go-test/cover.html",
		WorkDir: "/my-project",
	}).Return(nil, nil)

	action := TestAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}

func TestDebugTest(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On("ModuleAction", mock.Anything).Return(GoModTestData(true), nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "go test -vet off -cover -covermode=count -coverprofile /my-project/.dist/github-com-cidverse-my-project/go-test/cover.out -v ./...",
		WorkDir: "/my-project",
	}).Return(nil, nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command:       "go test -coverprofile /my-project/.dist/github-com-cidverse-my-project/go-test/cover.out -json -covermode=count ./...",
		WorkDir:       "/my-project",
		CaptureOutput: true,
	}).Return(&cidsdk.ExecuteCommandResponse{Stdout: "{}"}, nil)
	sdk.On("FileWrite", "/my-project/.dist/github-com-cidverse-my-project/go-test/cover.json", []byte("{}")).Return(nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "go tool cover -html /my-project/.dist/github-com-cidverse-my-project/go-test/cover.out -o /my-project/.dist/github-com-cidverse-my-project/go-test/cover.html",
		WorkDir: "/my-project",
	}).Return(nil, nil)

	action := TestAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}

func TestUnsupportedTest(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On("ModuleAction", mock.Anything).Return(api.GetUnknownTestData(false), nil)

	action := TestAction{Sdk: sdk}
	err := action.Execute()
	assert.Error(t, errors.New("build system unknown is not supported"), err)
}
