package applicationinspector

import (
	"testing"

	"github.com/cidverse/cid-actions-go/actions/api"
	"github.com/cidverse/cid-actions-go/pkg/core/test"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestApplicationInspectorScanBuild(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ProjectAction", mock.Anything).Return(api.GetProjectActionData(false), nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: `appinspector analyze --no-show-progress -s "/my-project" --base-path "/my-project" --repository-uri "https://github.com/cidverse/normalizeci.git" --commit-hash "abcdef123456" -f json -o "/my-project/.tmp/applicationinspector.json" -g **/tests/**,**/.git/**,**/.dist/**,**/.tmp/**`,
		WorkDir: `/my-project`,
	}).Return(nil, nil)
	sdk.On("ArtifactUpload", cidsdk.ArtifactUploadRequest{
		File:          "/my-project/.tmp/applicationinspector.json",
		Type:          "report",
		Format:        "applicationinspector",
		FormatVersion: "json",
	}).Return(nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
