package sonarqube

import (
	"testing"

	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/cidverse/cid-sdk-go/mocks"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSonarqubeScanGoMod(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On("ProjectAction", &ScanConfig{}).Return(SonarqubeGoModTestData(false), nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*ScanConfig)
		arg.SonarHostURL = "https://sonarcloud.local"
		arg.SonarOrganization = "my-org"
		arg.SonarProjectKey = "my-project-key"
		arg.SonarToken = "my-token"
	})
	sdk.On("ArtifactList", cidsdk.ArtifactListRequest{ArtifactType: "report", Format: "sarif"}).Return(&[]cidsdk.ActionArtifact{
		{
			BuildID:       "0",
			JobID:         "0",
			Module:        "root",
			Type:          "report",
			Name:          "test.sarif.json",
			Format:        "sarif",
			FormatVersion: "2.1.0",
		},
	}, nil)
	sdk.On("ArtifactDownload", cidsdk.ArtifactDownloadRequest{
		Module:     "root",
		Type:       "report",
		Name:       "test.sarif.json",
		TargetFile: "/my-project/.tmp/test.sarif.json",
	}).Return(nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "sonar-scanner -D sonar.host.url=https://sonarcloud.local -D sonar.login=my-token -D sonar.organization=my-org -D sonar.projectKey=my-project-key -D sonar.projectName=my-project-name -D sonar.branch.name= -D sonar.sources=. -D sonar.tests=. -D sonar.sarifReportPaths=/my-project/.tmp/test.sarif.json -D sonar.go.coverage.reportPaths=/my-project/.dist/github-com-cidverse-my-project/go-test/coverage.out -D sonar.go.tests.reportPaths=/my-project/.dist/github-com-cidverse-my-project/go-test/coverage.json -D sonar.inclusions= -D sonar.exclusions=**/*_test.go,**/vendor/**,**/testdata/* -D sonar.test.inclusions=**/*_test.go -D sonar.test.exclusions=**/vendor/**",
		WorkDir: "/my-project",
	}).Return(nil, nil)

	httpmock.ActivateNonDefault(apiClient.GetClient())
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "https://sonarcloud.local/api/projects/create?name=my-project-name&organization=my-org&project=my-project-key", httpmock.NewStringResponder(200, ``))
	httpmock.RegisterResponder("POST", "https://sonarcloud.local/api/project_branches/rename?name=develop&project=my-project-key", httpmock.NewStringResponder(200, ``))

	action := ScanAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
