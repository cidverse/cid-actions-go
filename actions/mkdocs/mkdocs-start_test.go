package mkdocs

import (
	"testing"

	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/cidverse/cid-sdk-go/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTechdocsStart(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On(`ModuleAction`, mock.Anything).Return(MKDocsTestData(false), nil)
	sdk.On(`ExecuteCommand`, cidsdk.ExecuteCommandRequest{
		Command: `pipenv sync`,
		WorkDir: `/my-project/docs`,
	}).Return(nil, nil)
	sdk.On(`ExecuteCommand`, cidsdk.ExecuteCommandRequest{
		Command: `pipenv run mkdocs serve --dev-addr 0.0.0.0:7600 --watch /my-project/docs`,
		WorkDir: `/my-project/docs`,
		Ports: []int{7600},
	}).Return(nil, nil)

	action := StartAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
