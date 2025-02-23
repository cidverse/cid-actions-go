package javabuild

import (
	"fmt"
	"strings"

	"github.com/cidverse/cid-actions-go/actions/java/javacommon"
	cidsdk "github.com/cidverse/cid-sdk-go"
)

type Action struct {
	Sdk cidsdk.SDKClient
}

type Config struct {
	MavenVersion string `json:"maven_version"        env:"MAVEN_VERSION"`
}

func (a Action) Metadata() cidsdk.ActionMetadata {
	return cidsdk.ActionMetadata{
		Name:        "java-build",
		Description: `Builds the java module using the configured build system.`,
		Category:    "build",
		Scope:       cidsdk.ActionScopeModule,
		Rules: []cidsdk.ActionRule{
			{
				Type:       "cel",
				Expression: `MODULE_BUILD_SYSTEM == "gradle"`,
			},
			{
				Type:       "cel",
				Expression: `MODULE_BUILD_SYSTEM == "maven"`,
			},
		},
		Access: cidsdk.ActionAccess{
			Environment: []cidsdk.ActionAccessEnv{},
			Executables: []cidsdk.ActionAccessExecutable{},
		},
	}
}

func (a Action) Execute() (err error) {
	cfg := Config{}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
	}

	// version
	if cfg.MavenVersion == "" {
		cfg.MavenVersion = javacommon.GetVersion(ctx.Env["NCI_COMMIT_REF_TYPE"], ctx.Env["NCI_COMMIT_REF_RELEASE"], ctx.Env["NCI_COMMIT_HASH_SHORT"])
	}

	// build
	if ctx.Module.BuildSystem == string(cidsdk.BuildSystemGradle) {
		// verify gradle wrapper
		err = javacommon.VerifyGradleWrapper(ctx.Module.ModuleDir)
		if err != nil {
			return err
		}

		gradleWrapper := cidsdk.JoinPath(ctx.Module.ModuleDir, "gradlew")
		if !a.Sdk.FileExists(gradleWrapper) {
			return fmt.Errorf("gradle wrapper not found at %s", gradleWrapper)
		}

		buildArgs := []string{
			fmt.Sprintf(`-Pversion=%q`, cfg.MavenVersion),
			`assemble`,
			`--no-daemon`,
			`--warning-mode=all`,
			`--console=plain`,
			`--stacktrace`,
		}
		buildResult, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: fmt.Sprintf("java-exec %s %s", gradleWrapper, strings.Join(buildArgs, " ")),
			WorkDir: ctx.Module.ModuleDir,
		})
		if err != nil {
			return err
		} else if buildResult.Code != 0 {
			return fmt.Errorf("gradle build failed, exit code %d", buildResult.Code)
		}
	} else if ctx.Module.BuildSystem == string(cidsdk.BuildSystemMaven) {
		mavenWrapper := cidsdk.JoinPath(ctx.Module.ModuleDir, "mvnw")
		if !a.Sdk.FileExists(mavenWrapper) {
			return fmt.Errorf("maven wrapper not found at %s", mavenWrapper)
		}

		buildArgs := []string{
			`package`,
			`--batch-mode`,
			`-Dmaven.test.skip=true`,
		}
		buildResult, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: fmt.Sprintf("java-exec %s %s", mavenWrapper, strings.Join(buildArgs, " ")),
			WorkDir: ctx.Module.ModuleDir,
		})
		if err != nil {
			return err
		} else if buildResult.Code != 0 {
			return fmt.Errorf("maven build failed, exit code %d", buildResult.Code)
		}
	}

	return nil
}
