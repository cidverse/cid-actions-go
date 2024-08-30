package ansiblelint

import (
	"testing"

	"github.com/cidverse/cid-actions-go/actions/api"
	"github.com/cidverse/cid-actions-go/pkg/core/test"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAnsibleLint(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ModuleAction", mock.Anything).Return(api.GetAnsibleTestData(false), nil)
	sdk.On("FileExists", "/my-project/playbook-a/requirements.yml").Return(false)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "ansible-lint .",
		WorkDir: "/my-project/playbook-a",
	}).Return(nil, nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}

func TestAnsibleLintWithDependencies(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ModuleAction", mock.Anything).Return(api.GetAnsibleTestData(false), nil)
	sdk.On("FileExists", "/my-project/playbook-a/requirements.yml").Return(true)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "ansible-galaxy collection install -r requirements.yml",
		WorkDir: "/my-project/playbook-a",
	}).Return(nil, nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "ansible-lint .",
		WorkDir: "/my-project/playbook-a",
	}).Return(nil, nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
