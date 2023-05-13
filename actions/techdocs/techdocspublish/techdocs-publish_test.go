package techdocspublish

import (
	"testing"

	"github.com/cidverse/cid-actions-go/actions/techdocs/techdocscommon"
	"github.com/cidverse/cid-actions-go/pkg/core/test"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTechdocsS3Publish(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On(`ModuleAction`, mock.Anything).Return(techdocscommon.TechdocsTestData(false), nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*Config)
		arg.Entity = `default/component/my-entity`
		arg.Target = `awsS3`
		arg.S3Endpoint = `minio.local`
		arg.S3Bucket = `techdocs`
		arg.S3AccessKey = `123456`
		arg.S3SecretKey = `123456abcdef`
		arg.S3ForcePathStyle = true
	})
	sdk.On("ArtifactDownload", cidsdk.ArtifactDownloadRequest{
		ID:         "my-module|html|docs.tar",
		TargetFile: "/my-project/.tmp/docs.tar",
	}).Return(nil)
	sdk.On("TARExtract", "/my-project/.tmp/docs.tar", "/my-project/.tmp/public").Return(nil)
	sdk.On(`ExecuteCommand`, cidsdk.ExecuteCommandRequest{
		Command: `techdocs-cli publish --entity default/component/my-entity --directory /my-project/.tmp/public --publisher-type awsS3 --awsEndpoint minio.local --storage-name techdocs --awsS3ForcePathStyle`,
		WorkDir: `/my-project`,
		Env: map[string]string{
			`AWS_ENDPOINT`:          `minio.local`,
			`AWS_ACCESS_KEY_ID`:     `123456`,
			`AWS_SECRET_ACCESS_KEY`: `123456abcdef`,
			`AWS_REGION`:            ``,
		},
	}).Return(&cidsdk.ExecuteCommandResponse{Code: 0}, nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
