package techdocs

import (
	"testing"

	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/cidverse/cid-sdk-go/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTechdocsS3Publish(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On(`ModuleAction`, mock.Anything).Return(TechdocsTestData(false), nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*PublishConfig)
		arg.Entity = `default/component/my-entity`
		arg.Target = `awsS3`
		arg.S3Endpoint = `minio.local`
		arg.S3Bucket = `techdocs`
		arg.S3AccessKey = `123456`
		arg.S3SecretKey = `123456abcdef`
		arg.S3ForcePathStyle = true
	})
	sdk.On(`ExecuteCommand`, cidsdk.ExecuteCommandRequest{
		Command: `techdocs-cli publish --entity default/component/my-entity --directory /my-project/.dist/my-module/html --publisher-type awsS3 --awsEndpoint minio.local --storage-name techdocs --awsS3ForcePathStyle`,
		WorkDir: `/my-project`,
		Env: map[string]string{
			`AWS_ACCESS_KEY_ID`:     `123456`,
			`AWS_SECRET_ACCESS_KEY`: `123456abcdef`,
		},
	}).Return(nil, nil)

	action := PublishAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
