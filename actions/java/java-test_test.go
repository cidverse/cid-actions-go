package java

import (
	"testing"

	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/cidverse/cid-sdk-go/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestJavaTestGradle(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On("ModuleAction", mock.Anything).Return(GradleTestData(false), nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: `java --add-opens=java.prefs/java.util.prefs=ALL-UNNAMED "-Dorg.gradle.appname=gradlew" "-classpath" "gradle/wrapper/gradle-wrapper.jar" "org.gradle.wrapper.GradleWrapperMain" -Pversion="main-SNAPSHOT" check --no-daemon --warning-mode=all --console=plain --stacktrace`,
		WorkDir: "/my-project",
	}).Return(&cidsdk.ExecuteCommandResponse{Code: 0}, nil)
	sdk.On("FileList", cidsdk.FileRequest{Directory: "/my-project", Extensions: []string{"jacocoTestReport.xml"}}).Return([]cidsdk.File{cidsdk.NewFile("/my-project/build/reports/jacoco/test/jacocoTestReport.xml")}, nil)
	sdk.On("ArtifactUpload", cidsdk.ArtifactUploadRequest{
		Module: "my-module",
		File:   "/my-project/build/reports/jacoco/test/jacocoTestReport.xml",
		Type:   "report",
		Format: "jacoco",
	}).Return(nil)

	action := TestAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
