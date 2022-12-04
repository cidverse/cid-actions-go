package container

import (
	"testing"

	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/cidverse/cid-sdk-go/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestContainerPublishOCIManifest(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On("ModuleAction", mock.Anything).Return(DockerfileTestData(false), nil).Run(func(args mock.Arguments) {
		// arg := args.Get(0).(*BuildConfig)
	})
	sdk.On("FileRead", "/my-project/.dist/my-project/oci-image/image.txt").Return("localhost/my-image:latest", nil)
	sdk.On("FileList", cidsdk.FileRequest{Directory: "/my-project/.dist/my-project/oci-image", Extensions: []string{".tar"}}).Return([]cidsdk.File{cidsdk.NewFile("/my-project/.dist/my-project/oci-image/linux_amd64.tar")}, nil)
	sdk.On("UUID").Return("58b1ee45-58a6-4c1d-be68-15aa1ce24268")
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "buildah manifest create 58b1ee4558a64c1dbe6815aa1ce24268",
		WorkDir: "/my-project",
	}).Return(nil, nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "buildah manifest add 58b1ee4558a64c1dbe6815aa1ce24268 oci-archive:/my-project/.dist/my-project/oci-image/linux_amd64.tar",
		WorkDir: "/my-project",
	}).Return(nil, nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "buildah manifest inspect 58b1ee4558a64c1dbe6815aa1ce24268",
		WorkDir: "/my-project",
	}).Return(nil, nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "buildah manifest push --all --format oci 58b1ee4558a64c1dbe6815aa1ce24268 docker://localhost/my-image:latest",
		WorkDir: "/my-project",
	}).Return(nil, nil)

	action := BuildahPublishAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
