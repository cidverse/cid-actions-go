package helm

import (
	"testing"

	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/cidverse/cid-sdk-go/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHelmLint(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On("ModuleAction", &LintConfig{}).Return(GetHelmTestData(false), nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "helm lint . --strict",
		WorkDir: "/my-project",
	}).Return(nil, nil)

	action := LintAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
