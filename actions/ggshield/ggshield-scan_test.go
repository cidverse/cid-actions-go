package ggshield

import (
	"testing"

	"github.com/cidverse/cid-actions-go/actions/api"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/cidverse/cid-sdk-go/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGGShieldScanBuild(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On("ProjectAction", &ScanConfig{}).Return(api.GetProjectActionData(false), nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "ggshield scan path -r -y .",
		WorkDir: "/my-project",
	}).Return(nil, nil)

	action := ScanAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
