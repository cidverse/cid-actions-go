package helm

/*
func TestHelmPublishNexus(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On("ModuleAction", &PublishNexusConfig{}).Return(GetHelmTestData(false), nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*PublishNexusConfig)
		arg.NexusURL = "localhost:9999"
		arg.NexusRepository = "dummy"
		arg.NexusUsername = "admin"
		arg.NexusPassword = "admin"
	})
	sdk.On("FileList", cidsdk.FileRequest{Directory: ".dist/my-package/helm-charts", Extensions: []string{".tgz"}}).Return([]cidsdk.File{cidsdk.NewFile(".dist/my-package/helm-charts/my-chart.tgz")}, nil)

	action := PublishNexusAction{Sdk: sdk}
	err := action.Execute()
	assert.Error(t, err)
}
*/
