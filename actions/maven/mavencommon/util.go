package mavencommon

import (
	"fmt"
)

const wrapperJar = ".mvn/wrapper/maven-wrapper.jar"

func MavenWrapperCommand(args string) string {
	return fmt.Sprintf("java -classpath=%q org.apache.maven.wrapper.MavenWrapperMain %s", wrapperJar, args)
}
