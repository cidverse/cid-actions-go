package semgrepscan

import (
	_ "embed"
	"testing"

	"github.com/cidverse/cid-actions-go/pkg/core/test"
	"github.com/cidverse/cid-actions-go/testdata"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//go:embed files/report.sarif.json
var reportJson string

func TestSemgrepScan(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ProjectAction", mock.Anything).Return(testdata.ModuleDefault(nil, false), nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command:       "semgrep ci --sarif --quiet --metrics=off --disable-version-check --exclude=.dist --exclude=.tmp --config \"p/ci\"",
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

	action := ScanAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}

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

	action := ScanAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
