package mavenbuild

import (
	"testing"

	"github.com/cidverse/cid-actions-go/actions/maven/mavencommon"
	"github.com/cidverse/cid-actions-go/pkg/core/test"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/stretchr/testify/assert"
)

func TestMavenBuild(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ModuleActionDataV1").Return(mavencommon.MavenTestData(map[string]string{
		"WRAPPER_VERIFICATION": "false",
	}, false), nil)
	sdk.On("FileExists", "/my-project/mvnw").Return(true)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: `java -classpath=".mvn/wrapper/maven-wrapper.jar" org.apache.maven.wrapper.MavenWrapperMain versions:set -DnewVersion="1.0.0"`,
		WorkDir: "/my-project",
	}).Return(&cidsdk.ExecuteCommandResponse{Code: 0}, nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: `java -classpath=".mvn/wrapper/maven-wrapper.jar" org.apache.maven.wrapper.MavenWrapperMain package --batch-mode -Dmaven.test.skip=true`,
		WorkDir: "/my-project",
	}).Return(&cidsdk.ExecuteCommandResponse{Code: 0}, nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}
