package upxoptimize

import (
	"testing"

	"github.com/cidverse/cid-actions-go/actions/api"
	"github.com/cidverse/cid-actions-go/pkg/core/test"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ProjectAction", mock.Anything).Return(api.GetProjectActionData(false), nil)
	sdk.On("FileList", cidsdk.FileRequest{Directory: "/my-project/bin"}).Return([]cidsdk.File{cidsdk.NewFile("/my-project/bin/linux_amd64")}, nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "upx --lzma /my-project/bin/linux_amd64",
		WorkDir: "/my-project",
	}).Return(nil, nil)

	action := OptimizeAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
