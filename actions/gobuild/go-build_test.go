package gobuild

import (
	"errors"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/cidverse/cid-sdk-go/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

/*
func TestLocal(t *testing.T) {
	os.Setenv("CID_API_ADDR", "http://localhost:7400")
	sdk, err := cidsdk.NewSDK()
	assert.NoError(t, err)

	action := Action{Sdk: sdk}
	err = action.Execute()
	assert.NoError(t, err)
}*/

func TestGoMod(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On("Log", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	sdk.On("PrepareAction", &Config{}).Return(cidsdk.ActionEnv{
		Module: cidsdk.ProjectModule{
			ProjectDir:        "/my-project",
			ModuleDir:         "/my-project",
			Discovery:         []string{"file~/my-project/go.mod"},
			Name:              "github.com/cidverse/my-project",
			Slug:              "github-com-cidverse-my-project",
			BuildSystem:       "gomod",
			BuildSystemSyntax: "default",
			Language:          &map[string]string{"go": "1.19.0"},
			Submodules:        nil,
		},
		Config: cidsdk.CurrentConfig{
			Debug:       false,
			Log:         map[string]string{},
			ArtifactDir: ".dist",
			TempDir:     ".tmp",
		},
	}, nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*Config)
		arg.Platform = []Platform{
			{
				Goos:   "linux",
				Goarch: "amd64",
			},
		}
	})

	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "go build -buildvcs=false -o .dist/github-com-cidverse-my-project/bin/linux_amd64 .",
		WorkDir: "/my-project",
		Env: map[string]string{
			"CGO_ENABLED": "false",
			"GOPROXY":     "https://goproxy.io,direct",
			"GOOS":        "linux",
			"GOARCH":      "amd64",
		},
	}).Return(nil, nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}

func TestUnsupported(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On("PrepareAction", &Config{}).Return(cidsdk.ActionEnv{
		Module: cidsdk.ProjectModule{
			ProjectDir:        "/my-project",
			ModuleDir:         "/my-project",
			Discovery:         []string{"file~/my-project/go.mod"},
			Name:              "github.com/cidverse/my-project",
			Slug:              "github-com-cidverse-my-project",
			BuildSystem:       "unknown",
			BuildSystemSyntax: "default",
			Language:          &map[string]string{"go": "1.19.0"},
			Submodules:        nil,
		},
		Config: cidsdk.CurrentConfig{
			Debug:       false,
			Log:         map[string]string{},
			ArtifactDir: ".dist",
			TempDir:     ".tmp",
		},
	}, nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*Config)
		arg.Platform = []Platform{
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
