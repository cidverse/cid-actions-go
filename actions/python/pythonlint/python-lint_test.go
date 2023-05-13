package pythonlint

import (
	"testing"

	"github.com/cidverse/cid-actions-go/pkg/core/test"
	"github.com/cidverse/cid-actions-go/testdata"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPythonLint(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ModuleAction", mock.Anything).Return(testdata.ModulePython(string(cidsdk.BuildSystemRequirementsTXT)), nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "flake8 .",
		WorkDir: "/my-project",
	}).Return(nil, nil)

	action := LintAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
