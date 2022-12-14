package container

import (
	"testing"

	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/cidverse/cid-sdk-go/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestContainerBuildDockerfile(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On("ModuleAction", mock.Anything).Return(DockerfileTestData(false), nil).Run(func(args mock.Arguments) {
		// arg := args.Get(0).(*BuildConfig)
	})
	sdk.On("FileRead", "/my-project/Dockerfile").Return("FROM alpine:latest", nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "buildah build --platform linux/amd64 -f Dockerfile -t oci-archive:.dist/my-project/oci-image/linux_amd64.tar --layers --squash --annotation \"org.opencontainers.image.source=${NCI_REPOSITORY_REMOTE}\" --annotation \"org.opencontainers.image.created=${TIMESTAMP_RFC3339}\" --annotation \"org.opencontainers.image.authors=\" --annotation \"org.opencontainers.image.title=my-project\" --annotation \"org.opencontainers.image.description=\" /my-project",
		WorkDir: "/my-project",
	}).Return(nil, nil)
	sdk.On("FileWrite", "/my-project/.dist/my-project/oci-image/image.txt", []byte("ghcr.io/cidverse/dummy:latest")).Return(nil)

	action := BuildahBuildAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
