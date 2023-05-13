package containerpublish

import (
	"testing"

	"github.com/cidverse/cid-actions-go/pkg/core/test"
	"github.com/cidverse/cid-actions-go/testdata"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestContainerPublishOCI(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ModuleAction", mock.Anything).Return(testdata.ModuleDockerfile(), nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*Config)
		arg.AlwaysPublishManifest = false
	})
	sdk.On("FileRead", "/my-project/.dist/my-project/oci-image/image.txt").Return("localhost/my-image:latest", nil)
	sdk.On("FileList", cidsdk.FileRequest{Directory: "/my-project/.dist/my-project/oci-image", Extensions: []string{".tar"}}).Return([]cidsdk.File{cidsdk.NewFile("/my-project/.dist/my-project/oci-image/linux_amd64.tar")}, nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "skopeo copy --digestfile /my-project/.tmp/digest.txt oci-archive:/my-project/.dist/my-project/oci-image/linux_amd64.tar docker://localhost/my-image:latest",
		WorkDir: "/my-project",
	}).Return(&cidsdk.ExecuteCommandResponse{Code: 0}, nil)
	sdk.On("ArtifactUpload", cidsdk.ArtifactUploadRequest{
		Module: "my-project",
		File:   "/my-project/.tmp/digest.txt",
		Type:   "oci-image",
		Format: "digest",
	}).Return(nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}

func TestContainerPublishOCIManifest(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ModuleAction", mock.Anything).Return(testdata.ModuleDockerfile(), nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*Config)
		arg.AlwaysPublishManifest = true
	})
	sdk.On("FileRead", "/my-project/.dist/my-project/oci-image/image.txt").Return("localhost/my-image:latest", nil)
	sdk.On("FileList", cidsdk.FileRequest{Directory: "/my-project/.dist/my-project/oci-image", Extensions: []string{".tar"}}).Return([]cidsdk.File{cidsdk.NewFile("/my-project/.dist/my-project/oci-image/linux_amd64.tar")}, nil)
	sdk.On("UUID").Return("58b1ee45-58a6-4c1d-be68-15aa1ce24268")
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "buildah manifest create 58b1ee4558a64c1dbe6815aa1ce24268",
		WorkDir: "/my-project",
	}).Return(&cidsdk.ExecuteCommandResponse{Code: 0}, nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "buildah manifest add 58b1ee4558a64c1dbe6815aa1ce24268 oci-archive:/my-project/.dist/my-project/oci-image/linux_amd64.tar",
		WorkDir: "/my-project",
	}).Return(&cidsdk.ExecuteCommandResponse{Code: 0}, nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: "buildah manifest push --all --format oci --digestfile /my-project/.tmp/digest.txt 58b1ee4558a64c1dbe6815aa1ce24268 docker://localhost/my-image:latest",
		WorkDir: "/my-project",
	}).Return(&cidsdk.ExecuteCommandResponse{Code: 0}, nil)
	sdk.On("ArtifactUpload", cidsdk.ArtifactUploadRequest{
		Module: "my-project",
		File:   "/my-project/.tmp/digest.txt",
		Type:   "oci-image",
		Format: "digest",
	}).Return(nil)
	sdk.On("FileRead", "/my-project/.tmp/digest.txt").Return("sha256:db3c9370a2728b36f7e3389d5adbafd7f1e608413dc213b9d4a76972a35ca015", nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command:       "buildah manifest inspect localhost/my-image@sha256:db3c9370a2728b36f7e3389d5adbafd7f1e608413dc213b9d4a76972a35ca015",
		WorkDir:       "/my-project",
		CaptureOutput: true,
	}).Return(&cidsdk.ExecuteCommandResponse{Code: 0, Stdout: `{"schemaVersion":2,"mediaType":"application/vnd.oci.image.index.v1+json","manifests":[{"mediaType":"application/vnd.oci.image.manifest.v1+json","digest":"sha256:7b7ef082c7f44942b78edfa80e0a2c9c05437504be44b372e22a883e1ef76d08","size":877,"platform":{"architecture":"amd64","os":"linux"}}]}`}, nil)
	sdk.On("FileWrite", "/my-project/.tmp/manifest.json", []byte(`{"schemaVersion":2,"mediaType":"application/vnd.oci.image.index.v1+json","manifests":[{"mediaType":"application/vnd.oci.image.manifest.v1+json","digest":"sha256:7b7ef082c7f44942b78edfa80e0a2c9c05437504be44b372e22a883e1ef76d08","size":877,"platform":{"architecture":"amd64","os":"linux"}}]}`)).Return(nil)
	sdk.On("ArtifactUpload", cidsdk.ArtifactUploadRequest{
		Module:        "my-project",
		File:          "/my-project/.tmp/manifest.json",
		Type:          "oci-image",
		Format:        "manifest",
		FormatVersion: "oci",
	}).Return(nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
