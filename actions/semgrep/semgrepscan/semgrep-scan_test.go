package semgrepscan

import (
	_ "embed"
	"testing"

	"github.com/cidverse/cid-actions-go/actions/api"
	"github.com/cidverse/cid-actions-go/pkg/core/test"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//go:embed files/report.sarif.json
var reportJson string

func TestSemgrepScan(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ProjectAction", mock.Anything).Return(api.GetProjectActionData(false), nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command:       "semgrep scan --config p/default --sarif --quiet --metrics=off --disable-version-check --exclude=.dist --exclude=.tmp",
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
