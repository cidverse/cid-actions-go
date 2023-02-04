package java

import (
	"testing"

	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/cidverse/cid-sdk-go/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestJavaPublishGradle(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On("ModuleAction", mock.Anything).Return(GradleTestData(false), nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*PublishConfig)
		arg.MavenRepositoryUrl = "http://localhost:9100/test"
		arg.MavenRepositoryUsername = "admin"
		arg.MavenRepositoryPassword = "secret"
	})
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: `java --add-opens=java.prefs/java.util.prefs=ALL-UNNAMED "-Dorg.gradle.appname=gradlew" "-classpath" "gradle/wrapper/gradle-wrapper.jar" "org.gradle.wrapper.GradleWrapperMain" -Pversion="main-SNAPSHOT" publish --no-daemon --warning-mode=all --console=plain --stacktrace`,
		WorkDir: "/my-project",
		Env: map[string]string{
			"MAVEN_REPO_URL":      "http://localhost:9100/test",
			"MAVEN_REPO_USERNAME": "admin",
			"MAVEN_REPO_PASSWORD": "secret",
		},
	}).Return(&cidsdk.ExecuteCommandResponse{Code: 0}, nil)

	action := PublishAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
