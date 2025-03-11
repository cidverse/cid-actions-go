package javapublish

import (
	"testing"

	"github.com/cidverse/cid-actions-go/actions/java/javacommon"
	"github.com/cidverse/cid-actions-go/pkg/core/test"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/stretchr/testify/assert"
)

func TestJavaPublishGradle(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ModuleActionDataV1").Return(javacommon.GradleTestData(map[string]string{
		"WRAPPER_VERIFICATION": "false",
		"MAVEN_REPO_URL":       "http://localhost:9100/test",
		"MAVEN_REPO_USERNAME":  "admin",
		"MAVEN_REPO_PASSWORD":  "secret",
	}, false), nil)
	sdk.On("FileExists", "/my-project/gradlew").Return(true)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: `java-exec /my-project/gradlew -Pversion="1.0.0" publish --no-daemon --warning-mode=all --console=plain --stacktrace`,
		WorkDir: "/my-project",
		Env: map[string]string{
			"MAVEN_REPO_URL":      "http://localhost:9100/test",
			"MAVEN_REPO_USERNAME": "admin",
			"MAVEN_REPO_PASSWORD": "secret",
		},
	}).Return(&cidsdk.ExecuteCommandResponse{Code: 0}, nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
