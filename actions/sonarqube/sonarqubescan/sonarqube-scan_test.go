package sonarqubescan

import (
	"os"
	"testing"

	"github.com/cidverse/cid-actions-go/actions/api"
	"github.com/cidverse/cid-actions-go/actions/sonarqube/sonarqubecommon"
	"github.com/cidverse/cid-actions-go/pkg/core/test"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSonarqubeScanGoMod(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ProjectAction", &Config{}).Return(api.SonarqubeGoModTestData(false), nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*Config)
		arg.SonarHostURL = "https://sonarcloud.local"
		arg.SonarOrganization = "my-org"
		arg.SonarProjectKey = "my-project-key"
		arg.SonarToken = "my-token"
	})
	sdk.On("ArtifactList", cidsdk.ArtifactListRequest{Query: `artifact_type == "report"`}).Return(&[]cidsdk.ActionArtifact{
		{
			BuildID:       "0",
			JobID:         "0",
			ID:            "root|report|test.sarif.json",
			Module:        "root",
			Type:          "report",
			Name:          "test.sarif.json",
			Format:        "sarif",
			FormatVersion: "2.1.0",
		},
		{
			BuildID:       "0",
			JobID:         "0",
			ID:            "root|report|coverage.out",
			Module:        "root",
			Type:          "report",
			Name:          "coverage.out",
			Format:        "go-coverage",
			FormatVersion: "out",
		},
		{
			BuildID:       "0",
			JobID:         "0",
			ID:            "root|report|coverage.json",
			Module:        "root",
			Type:          "report",
			Name:          "coverage.json",
			Format:        "go-coverage",
			FormatVersion: "json",
		},
	}, nil)
	sdk.On("ArtifactDownload", cidsdk.ArtifactDownloadRequest{
		ID:         "root|report|test.sarif.json",
		TargetFile: "/my-project/.tmp/root-test.sarif.json",
	}).Return(nil)
	sdk.On("ArtifactDownload", cidsdk.ArtifactDownloadRequest{
		ID:         "root|report|coverage.out",
		TargetFile: "/my-project/.tmp/root-coverage.out",
	}).Return(nil)
	sdk.On("ArtifactDownload", cidsdk.ArtifactDownloadRequest{
		ID:         "root|report|coverage.json",
		TargetFile: "/my-project/.tmp/root-coverage.json",
	}).Return(nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: `sonar-scanner -X -D sonar.host.url=https://sonarcloud.local -D sonar.login=my-token -D sonar.projectKey=my-project-key -D sonar.projectName=my-project-name -D sonar.sources=. -D sonar.organization=my-org -D sonar.sarifReportPaths=/my-project/.tmp/root-test.sarif.json -D sonar.go.coverage.reportPaths=/my-project/.tmp/root-coverage.out -D sonar.go.tests.reportPaths=/my-project/.tmp/root-coverage.json -D sonar.exclusions=**/.git/**,**/*_test.go,**/vendor/**,**/mocks/**,**/testdata/* -D sonar.test.inclusions=**/*_test.go -D sonar.test.exclusions=**/vendor/** -D sonar.branch.name="main"`,
		WorkDir: "/my-project",
		Env: map[string]string{
			"SONAR_SCANNER_OPTS": " ",
		},
	}).Return(&cidsdk.ExecuteCommandResponse{Code: 0}, nil)

	httpmock.ActivateNonDefault(sonarqubecommon.ApiClient.GetClient())
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "https://sonarcloud.local/api/projects/create?name=my-project-name&organization=my-org&project=my-project-key", httpmock.NewStringResponder(200, ``))
	httpmock.RegisterResponder("POST", "https://sonarcloud.local/api/project_branches/rename?name=main&project=my-project-key", httpmock.NewStringResponder(200, ``))

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
	os.Clearenv()
}
