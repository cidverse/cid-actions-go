package gradlecommon

import (
	"fmt"
	"path/filepath"
	"strings"
)

const wrapperJar = "gradle/wrapper/gradle-wrapper.jar"

func GradleWrapperCommand(args string, rootDir string) string {
	appName := "gradlew"
	return fmt.Sprintf("java -Dorg.gradle.appname=%q -classpath %q org.gradle.wrapper.GradleWrapperMain %s", appName, filepath.Join(rootDir, wrapperJar), args)
}

// GetVersion returns the suggested java artifact version
// Unless the reference is a git tag versions will get a -SNAPSHOT suffix
func GetVersion(refType string, refName string, shortHash string) string {
	if refType == "tag" {
		return refName
	}

	return strings.TrimLeft(refName, "v") + "-" + shortHash
}
