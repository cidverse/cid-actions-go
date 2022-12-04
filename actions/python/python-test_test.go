package python

import (
	"testing"

	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/cidverse/cid-sdk-go/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPythonTest(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On("ModuleAction", mock.Anything).Return(PythonTestData(string(cidsdk.BuildSystemRequirementsTXT), false), nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "pytest",
		WorkDir: "/my-project",
	}).Return(nil, nil)

	action := TestAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
