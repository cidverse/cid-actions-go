package node

import (
	"testing"

	"github.com/cidverse/cid-actions-go/actions/api"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/cidverse/cid-sdk-go/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNodeBuild(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On("ModuleAction", &BuildConfig{}).Return(api.GetNodeTestData(false), nil)
	sdk.On("FileRead", "/my-project/package.json").Return(`{"scripts": {"build": ""}}`, nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "yarn install",
		WorkDir: "/my-project",
	}).Return(nil, nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "yarn build",
		WorkDir: "/my-project",
	}).Return(nil, nil)

	action := BuildAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}

func TestNodeBuildNoScript(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On("ModuleAction", &TestConfig{}).Return(api.GetNodeTestData(false), nil)
	sdk.On("FileRead", "/my-project/package.json").Return(`{}`, nil)

	action := TestAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
