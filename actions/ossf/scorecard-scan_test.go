package ossf

import (
	_ "embed"
	"testing"

	"github.com/cidverse/cid-actions-go/actions/api"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/cidverse/cid-sdk-go/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//go:embed files/report.sarif.json
var reportJson string

func TestOSSFScorecardScan(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On("ProjectAction", mock.Anything).Return(api.GetProjectActionData(false), nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command:       `scorecard --repo "https://github.com/cidverse/normalizeci" --format json --commit "abcdef123456" --checks "Contributors,Dependency-Update-Tool,Maintained,Security-Policy,Fuzzing,Branch-Protection,CI-Tests,Signed-Releases,Binary-Artifacts,SAST,License,Pinned-Dependencies,CII-Best-Practices,Code-Review,Dangerous-Workflow,Packaging,Token-Permissions,Vulnerabilities"`,
		WorkDir:       "/my-project",
		Env:           map[string]string{},
		CaptureOutput: true,
	}).Return(&cidsdk.ExecuteCommandResponse{Code: 0, Stdout: reportJson}, nil)
	sdk.On("ArtifactUploadByteArray", cidsdk.ArtifactUploadByteArrayRequest{
		File:          "ossf-scorecard.sarif.json",
		Content:       []byte(reportJson),
		Type:          "report",
		Format:        "sarif",
		FormatVersion: "2.1.0",
	}).Return(nil)

	action := ScorecardScanAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
