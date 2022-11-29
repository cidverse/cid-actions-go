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

func TestGoModBuild(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On("ModuleAction", mock.Anything).Return(GoModTestData(false), nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*BuildConfig)
		arg.Platform = []Platform{
			{
				Goos:   "linux",
				Goarch: "amd64",
			},
		}
	})

	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "go build -buildvcs=false -o /my-project/.dist/github-com-cidverse-my-project/bin/linux_amd64 .",
		WorkDir: "/my-project",
		Env: map[string]string{
			"CGO_ENABLED": "false",
			"GOPROXY":     "https://goproxy.io,direct",
			"GOOS":        "linux",
			"GOARCH":      "amd64",
		},
	}).Return(nil, nil)

	action := BuildAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}

func TestUnsupportedBuild(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On("ModuleAction", mock.Anything).Return(api.GetUnknownTestData(false), nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*BuildConfig)
		arg.Platform = []Platform{
			{
				Goos:   "linux",
				Goarch: "amd64",
			},
		}
	})

	action := BuildAction{Sdk: sdk}
	err := action.Execute()
	assert.Error(t, errors.New("build system unknown is not supported"), err)
}
