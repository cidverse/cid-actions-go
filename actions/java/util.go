package java

const javaGradleCmd = `java --add-opens=java.prefs/java.util.prefs=ALL-UNNAMED "-Dorg.gradle.appname=gradlew" "-classpath" "gradle/wrapper/gradle-wrapper.jar" "org.gradle.wrapper.GradleWrapperMain"`

// GetVersion returns the suggested java artifact version
// Unless the reference is a git tag versions will get a -SNAPSHOT suffix
func GetVersion(refType string, refName string) string {
	if refType == "tag" {
		return refName
	}

	return refName + "-SNAPSHOT"
}
