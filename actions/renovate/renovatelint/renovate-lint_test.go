package renovatelint

import (
	"testing"

	"github.com/cidverse/cid-actions-go/pkg/core/test"
	"github.com/cidverse/cid-actions-go/testdata"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRenovateLint(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ModuleAction", mock.Anything).Return(testdata.ModuleRenovate(), nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "renovate-config-validator --strict",
		WorkDir: "/my-project",
	}).Return(nil, nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
