package containerbuild

import (
	"testing"

	"github.com/cidverse/cid-actions-go/pkg/core/test"
	"github.com/cidverse/cid-actions-go/testdata"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestContainerBuildDockerfile(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ModuleAction", mock.Anything).Return(testdata.ModuleDockerfile(), nil).Run(func(args mock.Arguments) {
		// arg := args.Get(0).(*BuildConfig)
	})
	sdk.On("FileRead", "/my-project/Dockerfile").Return("FROM alpine:latest", nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "buildah build --platform linux/amd64 -f Dockerfile -t oci-archive:.dist/my-project/oci-image/linux_amd64.tar --layers --squash --annotation \"org.opencontainers.image.source={NCI_REPOSITORY_REMOTE}\" --annotation \"org.opencontainers.image.created={TIMESTAMP_RFC3339}\" --annotation \"org.opencontainers.image.authors=\" --annotation \"org.opencontainers.image.title=my-project\" --annotation \"org.opencontainers.image.description=\" /my-project",
		WorkDir: "/my-project",
	}).Return(&cidsdk.ExecuteCommandResponse{Code: 0}, nil)
	sdk.On("ArtifactUpload", cidsdk.ArtifactUploadRequest{
		File:          "image.txt",
		Content:       "ghcr.io/cidverse/dummy:latest",
		Module:        "my-project",
		Type:          "oci-image",
		Format:        "container-ref",
		FormatVersion: "",
	}).Return(nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
