package python

import (
	"testing"

	"github.com/cidverse/cid-actions-go/actions/api"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/cidverse/cid-sdk-go/mocks"
	"github.com/stretchr/testify/assert"
)

func TestPythonLint(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On("ModuleAction", &LintConfig{}).Return(api.GetPythonTestData(string(cidsdk.BuildSystemRequirementsTXT), false), nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "flake8 .",
		WorkDir: "/my-project",
	}).Return(nil, nil)

	action := LintAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
