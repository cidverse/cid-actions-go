package slsa

import (
	"testing"

	"github.com/cidverse/cid-sdk-go/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSLSAGenerate(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On(`ModuleAction`, mock.Anything).Return(SLSATestData(false), nil)

	action := GenerateAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
