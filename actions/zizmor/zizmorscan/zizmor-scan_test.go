package zizmorscan

import (
	"testing"

	"github.com/cidverse/cid-actions-go/pkg/core/test"
	"github.com/cidverse/cid-actions-go/testdata"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/stretchr/testify/assert"
)

func TestZizmorScan(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ProjectActionDataV1").Return(testdata.ModuleDefault(nil, false), nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: `zizmor . --format sarif --no-online-audits`,
		WorkDir: "/my-project",
		Env: map[string]string{
			"GH_HOST":  "github.com",
			"GH_TOKEN": "",
		},
		CaptureOutput: true,
	}).Return(&cidsdk.ExecuteCommandResponse{
		Code: 0,
	}, nil)
	sdk.On("FileWrite", "/my-project/.tmp/zizmor.sarif.json", []byte{}).Return(nil)
	sdk.On("ArtifactUpload", cidsdk.ArtifactUploadRequest{
		File:          "/my-project/.tmp/zizmor.sarif.json",
		Type:          "report",
		Format:        "sarif",
		FormatVersion: "2.1.0",
	}).Return(nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}

func TestZizmorScanIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
}
