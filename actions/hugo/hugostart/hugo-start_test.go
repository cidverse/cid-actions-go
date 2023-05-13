package hugostart

import (
	"testing"

	"github.com/cidverse/cid-actions-go/actions/hugo/hugocommon"
	"github.com/cidverse/cid-actions-go/pkg/core/test"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTechdocsStart(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On(`ModuleAction`, mock.Anything).Return(hugocommon.HugoTestData(false), nil)
	sdk.On(`ExecuteCommand`, cidsdk.ExecuteCommandRequest{
		Command: `hugo server --source /my-project/docs --minify --gc --baseUrl / --watch --port 7600`,
		WorkDir: `/my-project/docs`,
		Ports:   []int{7600},
	}).Return(nil, nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
