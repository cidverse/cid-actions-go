package helmlint

import (
	"testing"

	"github.com/cidverse/cid-actions-go/actions/helm/helmcommon"
	"github.com/cidverse/cid-actions-go/pkg/core/test"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/stretchr/testify/assert"
)

func TestHelmLint(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ModuleAction", &LintConfig{}).Return(helmcommon.GetHelmTestData(false), nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "helm lint . --strict",
		WorkDir: "/my-project/charts/mychart",
	}).Return(nil, nil)

	action := LintAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
