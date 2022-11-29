package fossa

import (
	"testing"

	"github.com/cidverse/cid-actions-go/actions/api"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/cidverse/cid-sdk-go/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGGShieldScanBuild(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On("ProjectAction", &SourceScanConfig{}).Return(api.GetProjectActionData(false), nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "fossa analyze",
		WorkDir: "/my-project",
	}).Return(nil, nil)

	action := SourceScanAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
