package qodana

import (
	_ "embed"
	"testing"

	"github.com/cidverse/cid-actions-go/actions/api"
	"github.com/cidverse/cid-sdk-go/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestQodanaScan(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On("ProjectAction", mock.Anything).Return(api.GetProjectActionData(false), nil)

	action := ScanAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
