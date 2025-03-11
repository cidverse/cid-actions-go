package renovatelint

import (
	"testing"

	"github.com/cidverse/cid-actions-go/pkg/core/test"
	"github.com/cidverse/cid-actions-go/testdata"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/stretchr/testify/assert"
)

func TestRenovateLint(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ProjectActionDataV1").Return(testdata.ProjectRenovate(), nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "renovate-config-validator --strict .",
		WorkDir: "/my-project",
	}).Return(&cidsdk.ExecuteCommandResponse{Code: 0}, nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
