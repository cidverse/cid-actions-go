package hugo

import (
	"testing"

	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/cidverse/cid-sdk-go/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTechdocsStart(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On(`ModuleAction`, mock.Anything).Return(HugoTestData(false), nil)
	sdk.On(`ExecuteCommand`, cidsdk.ExecuteCommandRequest{
		Command: `hugo server --source /my-project/docs --minify --gc --baseUrl / --watch --port 7600`,
		WorkDir: `/my-project/docs`,
		Ports: []int{7600},
	}).Return(nil, nil)

	action := StartAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
