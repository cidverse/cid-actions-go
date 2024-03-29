package nodetest

import (
	"testing"

	"github.com/cidverse/cid-actions-go/actions/api"
	"github.com/cidverse/cid-actions-go/pkg/core/test"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNodeTest(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ModuleAction", mock.Anything).Return(api.GetNodeTestData(false), nil)
	sdk.On("FileRead", "/my-project/package.json").Return(`{"scripts": {"test": ""}}`, nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "yarn install",
		WorkDir: "/my-project",
	}).Return(nil, nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "yarn test",
		WorkDir: "/my-project",
	}).Return(nil, nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}

func TestNodeTestNoScript(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ModuleAction", mock.Anything).Return(api.GetNodeTestData(false), nil)
	sdk.On("FileRead", "/my-project/package.json").Return(`{}`, nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
