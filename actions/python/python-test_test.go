package python

import (
	"testing"

	"github.com/cidverse/cid-actions-go/actions/api"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/cidverse/cid-sdk-go/mocks"
	"github.com/stretchr/testify/assert"
)

func TestPythonTest(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On("ModuleAction", &TestConfig{}).Return(api.GetPythonTestData(string(cidsdk.BuildSystemRequirementsTXT), false), nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "pytest",
		WorkDir: "/my-project",
	}).Return(nil, nil)

	action := TestAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
