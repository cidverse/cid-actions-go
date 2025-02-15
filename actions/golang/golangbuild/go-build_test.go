package golangbuild

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

func TestGoModBuild(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ModuleAction", mock.Anything).Return(golangcommon.GoModTestData(false), nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*BuildConfig)
		arg.Platform = []golangcommon.Platform{
			{
				Goos:   "linux",
				Goarch: "amd64",
			},
		}
	})

	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: `go get -v -t ./...`,
		WorkDir: "/my-project",
		Env: map[string]string{
			"GOTOOLCHAIN": "local",
		},
	}).Return(&cidsdk.ExecuteCommandResponse{Code: 0}, nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: `go build -buildvcs=false -ldflags "-s -w -X main.version={NCI_COMMIT_REF_RELEASE} -X main.commit={NCI_COMMIT_HASH} -X main.date={TIMESTAMP_RFC3339} -X main.status={NCI_REPOSITORY_STATUS}" -o /my-project/.tmp/linux_amd64 .`,
		WorkDir: "/my-project",
		Env: map[string]string{
			"CGO_ENABLED": "false",
			"GOOS":        "linux",
			"GOARCH":      "amd64",
			"GOTOOLCHAIN": "local",
		},
	}).Return(&cidsdk.ExecuteCommandResponse{Code: 0}, nil)
	sdk.On("ArtifactUpload", cidsdk.ArtifactUploadRequest{
		Module:        "github-com-cidverse-my-project",
		File:          "/my-project/.tmp/linux_amd64",
		Type:          "binary",
		Format:        "go",
		FormatVersion: "",
	}).Return(nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}

func TestUnsupportedBuild(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ModuleAction", mock.Anything).Return(api.GetUnknownTestData(false), nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*BuildConfig)
		arg.Platform = []golangcommon.Platform{
			{
				Goos:   "linux",
				Goarch: "amd64",
			},
		}
	})

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.Error(t, errors.New("build system unknown is not supported"), err)
}
