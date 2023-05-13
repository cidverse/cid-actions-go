package mkdocsstart

import (
	"testing"

	"github.com/cidverse/cid-actions-go/actions/mkdocs/mkdocscommon"
	"github.com/cidverse/cid-actions-go/pkg/core/test"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTechdocsStart(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On(`ModuleAction`, mock.Anything).Return(mkdocscommon.MKDocsTestData(false), nil)
	sdk.On(`ExecuteCommand`, cidsdk.ExecuteCommandRequest{
		Command: `pipenv sync`,
		WorkDir: `/my-project/docs`,
	}).Return(nil, nil)
	sdk.On(`ExecuteCommand`, cidsdk.ExecuteCommandRequest{
		Command: `pipenv run mkdocs serve --dev-addr 0.0.0.0:7600 --watch /my-project/docs`,
		WorkDir: `/my-project/docs`,
		Ports:   []int{7600},
	}).Return(nil, nil)

	action := StartAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
