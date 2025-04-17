package trivyfsscan

import (
	"testing"

	"github.com/cidverse/cid-actions-go/pkg/core/test"
	"github.com/cidverse/cid-actions-go/testdata"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/stretchr/testify/assert"
)

func TestTrivyFSScan(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ProjectActionDataV1").Return(testdata.ModuleDefault(nil, false), nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: `trivy fs . --severity MEDIUM,HIGH,CRITICAL --format sarif --output /my-project/.tmp/trivyfs.sarif.json`,
		WorkDir: "/my-project",
	}).Return(&cidsdk.ExecuteCommandResponse{
		Code: 0,
	}, nil)
	sdk.On("ArtifactUpload", cidsdk.ArtifactUploadRequest{
		File:          "/my-project/.tmp/trivyfs.sarif.json",
		Type:          "report",
		Format:        "sarif",
		FormatVersion: "2.1.0",
	}).Return(nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
