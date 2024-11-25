package ansiblelint

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

func TestAnsibleLint(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ModuleAction", mock.Anything).Return(api.GetAnsibleTestData(false), nil)
	sdk.On("FileExists", "/my-project/playbook-a/requirements.yml").Return(false)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "ansible-lint --project . --sarif-file /my-project/.tmp/ansiblelint.sarif.json",
		WorkDir: "/my-project/playbook-a",
	}).Return(nil, nil)
	sdk.On("FileRead", "/my-project/.tmp/ansiblelint.sarif.json").Return(reportJson, nil)
	sdk.On("ArtifactUpload", cidsdk.ArtifactUploadRequest{
		File:          "/my-project/.tmp/ansiblelint.sarif.json",
		Type:          "report",
		Format:        "sarif",
		FormatVersion: "2.1.0",
	}).Return(nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}

func TestAnsibleLintWithDependencies(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ModuleAction", mock.Anything).Return(api.GetAnsibleTestData(false), nil)
	sdk.On("FileExists", "/my-project/playbook-a/requirements.yml").Return(true)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "ansible-galaxy collection install -r requirements.yml",
		WorkDir: "/my-project/playbook-a",
	}).Return(nil, nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "ansible-lint --project . --sarif-file /my-project/.tmp/ansiblelint.sarif.json",
		WorkDir: "/my-project/playbook-a",
	}).Return(nil, nil)
	sdk.On("FileRead", "/my-project/.tmp/ansiblelint.sarif.json").Return(reportJson, nil)
	sdk.On("ArtifactUpload", cidsdk.ArtifactUploadRequest{
		File:          "/my-project/.tmp/ansiblelint.sarif.json",
		Type:          "report",
		Format:        "sarif",
		FormatVersion: "2.1.0",
	}).Return(nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
