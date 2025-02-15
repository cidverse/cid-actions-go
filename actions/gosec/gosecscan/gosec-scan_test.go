package gosecscan

import (
	_ "embed"
	"testing"

	"github.com/cidverse/cid-actions-go/actions/gosec/goseccommon"
	"github.com/cidverse/cid-actions-go/pkg/core/test"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//go:embed report.sarif.json
var reportJson string

func TestGosecScan(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ModuleAction", mock.Anything).Return(goseccommon.GoModTestData(false), nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "gosec -no-fail -fmt sarif -out /my-project/.tmp/gosec.sarif.json ./...",
		WorkDir: "/my-project",
	}).Return(nil, nil)
	sdk.On("FileRead", "/my-project/.tmp/gosec.sarif.json").Return(reportJson, nil)
	sdk.On("ArtifactUpload", cidsdk.ArtifactUploadRequest{
		File:          "/my-project/.tmp/gosec.sarif.json",
		Type:          "report",
		Format:        "sarif",
		FormatVersion: "2.1.0",
	}).Return(nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
