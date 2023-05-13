package hugobuild

import (
	"testing"

	"github.com/cidverse/cid-actions-go/actions/hugo/hugocommon"
	"github.com/cidverse/cid-actions-go/pkg/core/test"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHugoBuild(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On(`ModuleAction`, mock.Anything).Return(hugocommon.HugoTestData(false), nil)
	sdk.On(`ExecuteCommand`, cidsdk.ExecuteCommandRequest{
		Command: `hugo --minify --gc --source /my-project/docs --destination /my-project/.dist/my-module/html`,
		WorkDir: `/my-project`,
	}).Return(nil, nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
