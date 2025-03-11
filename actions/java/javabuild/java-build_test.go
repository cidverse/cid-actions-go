package javabuild

import (
	"testing"

	"github.com/cidverse/cid-actions-go/actions/java/javacommon"
	"github.com/cidverse/cid-actions-go/pkg/core/test"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/stretchr/testify/assert"
)

func TestJavaBuildGradle(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ModuleActionDataV1").Return(javacommon.GradleTestData(map[string]string{
		"WRAPPER_VERIFICATION": "false",
	}, false), nil)
	sdk.On("FileExists", "/my-project/gradlew").Return(true)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: `java-exec /my-project/gradlew -Pversion="1.0.0" assemble --no-daemon --warning-mode=all --console=plain --stacktrace`,
		WorkDir: "/my-project",
	}).Return(&cidsdk.ExecuteCommandResponse{Code: 0}, nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
