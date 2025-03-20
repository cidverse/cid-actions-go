package helmfilelint

import (
	"testing"

	"github.com/cidverse/cid-actions-go/actions/helm/helmcommon"
	"github.com/cidverse/cid-actions-go/pkg/core/test"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/stretchr/testify/assert"
)

func TestHelmfileLint(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ModuleActionDataV1").Return(helmcommon.GetHelmfileTestData(map[string]string{}, false), nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "helmfile lint",
		WorkDir: "/my-project",
	}).Return(&cidsdk.ExecuteCommandResponse{Code: 0}, nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
