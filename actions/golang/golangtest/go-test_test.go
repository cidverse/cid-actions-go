package golangtest

import (
	"errors"
	"testing"

	"github.com/cidverse/cid-actions-go/actions/api"
	"github.com/cidverse/cid-actions-go/actions/golang/golangcommon"
	"github.com/cidverse/cid-actions-go/pkg/core/test"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGoModTest(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ModuleAction", mock.Anything).Return(golangcommon.GoModTestData(false), nil)

	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: `go get -v -t ./...`,
		WorkDir: "/my-project",
		Env: map[string]string{
			"GOTOOLCHAIN": "local",
		},
	}).Return(&cidsdk.ExecuteCommandResponse{Code: 0}, nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "go test -vet off -cover -covermode=count -coverprofile /my-project/.tmp/cover.out -parallel=4 -timeout 10s ./...",
		WorkDir: "/my-project",
		Env: map[string]string{
			"GOTOOLCHAIN": "local",
		},
	}).Return(&cidsdk.ExecuteCommandResponse{Code: 0}, nil)
	sdk.On("ArtifactUpload", cidsdk.ArtifactUploadRequest{
		File:          "/my-project/.tmp/cover.out",
		Module:        "github-com-cidverse-my-project",
		Type:          "report",
		Format:        "go-coverage",
		FormatVersion: "out",
	}).Return(nil)

	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command:       "go test -coverprofile /my-project/.tmp/cover.out -json -covermode=count -parallel=4 -timeout 10s ./...",
		WorkDir:       "/my-project",
		CaptureOutput: true,
	}).Return(&cidsdk.ExecuteCommandResponse{Code: 0, Stdout: "{}"}, nil)
	sdk.On("FileWrite", "/my-project/.tmp/cover.json", []byte("{}")).Return(nil)
	sdk.On("ArtifactUpload", cidsdk.ArtifactUploadRequest{
		File:          "/my-project/.tmp/cover.json",
		Module:        "github-com-cidverse-my-project",
		Type:          "report",
		Format:        "go-coverage",
		FormatVersion: "json",
	}).Return(nil)

	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "go tool cover -html /my-project/.tmp/cover.out -o /my-project/.tmp/cover.html",
		WorkDir: "/my-project",
	}).Return(&cidsdk.ExecuteCommandResponse{Code: 0}, nil)
	sdk.On("ArtifactUpload", cidsdk.ArtifactUploadRequest{
		File:          "/my-project/.tmp/cover.html",
		Module:        "github-com-cidverse-my-project",
		Type:          "report",
		Format:        "go-coverage",
		FormatVersion: "html",
	}).Return(nil)

	action := TestAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}

func TestDebugTest(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ModuleAction", mock.Anything).Return(golangcommon.GoModTestData(true), nil)

	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: `go get -v -t ./...`,
		WorkDir: "/my-project",
		Env: map[string]string{
			"GOTOOLCHAIN": "local",
		},
	}).Return(&cidsdk.ExecuteCommandResponse{Code: 0}, nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "go test -vet off -cover -covermode=count -coverprofile /my-project/.tmp/cover.out -parallel=4 -timeout 10s -v ./...",
		WorkDir: "/my-project",
		Env: map[string]string{
			"GOTOOLCHAIN": "local",
		},
	}).Return(&cidsdk.ExecuteCommandResponse{Code: 0}, nil)
	sdk.On("ArtifactUpload", cidsdk.ArtifactUploadRequest{
		File:          "/my-project/.tmp/cover.out",
		Module:        "github-com-cidverse-my-project",
		Type:          "report",
		Format:        "go-coverage",
		FormatVersion: "out",
	}).Return(nil)

	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command:       "go test -coverprofile /my-project/.tmp/cover.out -json -covermode=count -parallel=4 -timeout 10s ./...",
		WorkDir:       "/my-project",
		CaptureOutput: true,
	}).Return(&cidsdk.ExecuteCommandResponse{Code: 0, Stdout: "{}"}, nil)
	sdk.On("FileWrite", "/my-project/.tmp/cover.json", []byte("{}")).Return(nil)
	sdk.On("ArtifactUpload", cidsdk.ArtifactUploadRequest{
		File:          "/my-project/.tmp/cover.json",
		Module:        "github-com-cidverse-my-project",
		Type:          "report",
		Format:        "go-coverage",
		FormatVersion: "json",
	}).Return(nil)

	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "go tool cover -html /my-project/.tmp/cover.out -o /my-project/.tmp/cover.html",
		WorkDir: "/my-project",
	}).Return(&cidsdk.ExecuteCommandResponse{Code: 0}, nil)
	sdk.On("ArtifactUpload", cidsdk.ArtifactUploadRequest{
		File:          "/my-project/.tmp/cover.html",
		Module:        "github-com-cidverse-my-project",
		Type:          "report",
		Format:        "go-coverage",
		FormatVersion: "html",
	}).Return(nil)

	action := TestAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}

func TestUnsupportedTest(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ModuleAction", mock.Anything).Return(api.GetUnknownTestData(false), nil)

	action := TestAction{Sdk: sdk}
	err := action.Execute()
	assert.Error(t, errors.New("build system unknown is not supported"), err)
}
