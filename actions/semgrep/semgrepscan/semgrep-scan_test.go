package semgrepscan

import (
	_ "embed"
	"testing"

	"github.com/cidverse/cid-actions-go/pkg/core/test"
	"github.com/cidverse/cid-actions-go/testdata"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/stretchr/testify/assert"
)

func TestSemgrepScan(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ProjectActionDataV1").Return(testdata.ModuleDefault(nil, false), nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: `semgrep ci --text --sarif-output="/my-project/.tmp/semgrep.sarif.json" --metrics=off --disable-version-check --exclude=.dist --exclude=.tmp --config "p/ci"`,
		WorkDir: "/my-project",
		Env: map[string]string{
			"SEMGREP_APP_TOKEN": "",
			"SEMGREP_RULES":     "",
		},
	}).Return(&cidsdk.ExecuteCommandResponse{
		Code: 0,
	}, nil)
	sdk.On("ArtifactUpload", cidsdk.ArtifactUploadRequest{
		File:          "/my-project/.tmp/semgrep.sarif.json",
		Type:          "report",
		Format:        "sarif",
		FormatVersion: "2.1.0",
	}).Return(nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}

/*
func TestSemgrepPRScan(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ProjectAction", mock.Anything).Return(testdata.ModuleDefault(map[string]string{
		"NCI_MERGE_REQUEST_SOURCE_HASH": "abcdef123456",
	}, false), nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command:       "semgrep ci --sarif --quiet --metrics=off --disable-version-check --exclude=.dist --exclude=.tmp --baseline abcdef123456 --config \"p/ci\"",
		WorkDir:       "/my-project",
		CaptureOutput: true,
	}).Return(&cidsdk.ExecuteCommandResponse{
		Code:   0,
		Stdout: reportJson,
	}, nil)
	sdk.On("FileWrite", "/my-project/.tmp/semgrep.sarif.json", []byte(reportJson)).Return(nil)
	sdk.On("ArtifactUpload", cidsdk.ArtifactUploadRequest{
		File:          "/my-project/.tmp/semgrep.sarif.json",
		Type:          "report",
		Format:        "sarif",
		FormatVersion: "2.1.0",
	}).Return(nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
*/
