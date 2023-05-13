package qodanascan

import (
	_ "embed"
	"testing"

	"github.com/cidverse/cid-actions-go/actions/qodana/qodanascancommon"
	"github.com/cidverse/cid-actions-go/pkg/core/test"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//go:embed qodana.sarif.json
var reportJson string

func TestQodanaScan(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ModuleAction", mock.Anything).Return(qodanascancommon.GoModuleTestData(), nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "qodana-go --source-directory=/my-project --results-dir=/my-project/.tmp --fail-threshold 10000",
		WorkDir: "/my-project",
		Env: map[string]string{
			"QODANA_BRANCH":     "",
			"QODANA_REMOTE_URL": "",
			"QODANA_REVISION":   "",
			"QODANA_TOKEN":      "",
		},
	}).Return(&cidsdk.ExecuteCommandResponse{Code: 0}, nil)
	sdk.On("FileRead", "/my-project/.tmp/qodana.sarif.json").Return(reportJson, nil)
	sdk.On("ArtifactUpload", cidsdk.ArtifactUploadRequest{
		Module:        "github-com-cidverse-my-project",
		File:          "/my-project/.tmp/qodana.sarif.json",
		Type:          "report",
		Format:        "sarif",
		FormatVersion: "2.1.0",
	}).Return(nil)

	action := ScanAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
