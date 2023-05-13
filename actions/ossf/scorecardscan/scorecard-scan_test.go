package scorecardscan

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

func TestOSSFScorecardScan(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ProjectAction", mock.Anything).Return(api.GetProjectActionData(false), nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command:       `scorecard --repo "https://github.com/cidverse/normalizeci" --format json --commit "abcdef123456" --checks "Contributors,Dependency-Update-Tool,Maintained,Security-Policy,Fuzzing,Branch-Protection,CI-Tests,Signed-Releases,Binary-Artifacts,SAST,License,Pinned-Dependencies,CII-Best-Practices,Code-Review,Dangerous-Workflow,Packaging,Token-Permissions,Vulnerabilities"`,
		WorkDir:       "/my-project",
		Env:           map[string]string{},
		CaptureOutput: true,
	}).Return(&cidsdk.ExecuteCommandResponse{Code: 0, Stdout: reportJson}, nil)
	sdk.On("ArtifactUpload", cidsdk.ArtifactUploadRequest{
		File:          "ossf-scorecard.sarif.json",
		Content:       reportJson,
		Type:          "report",
		Format:        "sarif",
		FormatVersion: "2.1.0",
	}).Return(nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
