package gitleaksscan

import (
	_ "embed"
	"testing"

	"github.com/cidverse/cid-actions-go/actions/api"
	"github.com/cidverse/cid-actions-go/pkg/core/test"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//go:embed report.sarif.json
var reportJson string

func TestGitleaksScanBuild(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ProjectAction", mock.Anything).Return(api.GetProjectActionData(false), nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: `gitleaks detect --source=. -v --no-git --report-format=sarif --report-path="/my-project/.tmp/gitleaks.sarif.json" --no-banner --redact=85 --exit-code 0`,
		WorkDir: "/my-project",
	}).Return(&cidsdk.ExecuteCommandResponse{Code: 0}, nil)
	sdk.On("FileRead", "/my-project/.tmp/gitleaks.sarif.json").Return(reportJson, nil)
	sdk.On("ArtifactUpload", cidsdk.ArtifactUploadRequest{
		File:          "/my-project/.tmp/gitleaks.sarif.json",
		Type:          "report",
		Format:        "sarif",
		FormatVersion: "2.1.0",
	}).Return(nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
