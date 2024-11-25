package ansibledeploy

import (
	"testing"

	"github.com/cidverse/cid-actions-go/actions/api"
	"github.com/cidverse/cid-actions-go/pkg/core/test"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAnsibleDeploy(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ModuleAction", mock.Anything).Return(api.GetAnsibleTestData(false), nil)
	sdk.On("FileExists", "/my-project/playbook-a/requirements.yml").Return(false)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: `ansible-playbook "/my-project/playbook-a/playbook.yml" -i "/my-project/playbook-a/inventory"`,
		WorkDir: "/my-project/playbook-a",
	}).Return(nil, nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
