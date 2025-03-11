package javatest

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
	MavenVersion        string `json:"maven_version"        env:"MAVEN_VERSION"`
	WrapperVerification bool   `json:"wrapper_verification" env:"WRAPPER_VERIFICATION"`
}

func (a Action) Metadata() cidsdk.ActionMetadata {
	return cidsdk.ActionMetadata{
		Name:        "java-test",
		Description: `Tests the java module using the configured build system.`,
		Category:    "test",
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
	// query action data
	d, err := a.Sdk.ModuleActionDataV1()
	if err != nil {
		return err
	}

	// parse config
	cfg := Config{}
	cidsdk.PopulateFromEnv(&cfg, d.Env)

	// version
	if cfg.MavenVersion == "" {
		cfg.MavenVersion = javacommon.GetVersion(d.Env["NCI_COMMIT_REF_TYPE"], d.Env["NCI_COMMIT_REF_RELEASE"], d.Env["NCI_COMMIT_HASH_SHORT"])
	}

	// test
	if d.Module.BuildSystem == string(cidsdk.BuildSystemGradle) {
		// verify gradle wrapper
		if cfg.WrapperVerification {
			err = javacommon.VerifyGradleWrapper(d.Module.ModuleDir)
			if err != nil {
				return err
			}
		}

		gradleWrapper := cidsdk.JoinPath(d.Module.ModuleDir, "gradlew")
		if !a.Sdk.FileExists(gradleWrapper) {
			return fmt.Errorf("gradle wrapper not found at %s", gradleWrapper)
		}

		testArgs := []string{
			fmt.Sprintf(`-Pversion=%q`, cfg.MavenVersion),
			`check`,
			`--no-daemon`,
			`--warning-mode=all`,
			`--console=plain`,
			`--stacktrace`,
		}
		testResult, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: fmt.Sprintf("java-exec %s %s", gradleWrapper, strings.Join(testArgs, " ")),
			WorkDir: d.Module.ModuleDir,
		})
		if err != nil {
			return err
		} else if testResult.Code != 0 {
			return fmt.Errorf("gradle test failed, exit code %d", testResult.Code)
		}

		// collect and store jacoco test reports
		testReports, err := a.Sdk.FileList(cidsdk.FileRequest{
			Directory:  d.Module.ModuleDir,
			Extensions: []string{"jacocoTestReport.xml"},
		})
		for _, report := range testReports {
			if strings.HasSuffix(report.Path, cidsdk.JoinPath("build", "reports", "jacoco", "test", "jacocoTestReport.xml")) {
				err := a.Sdk.ArtifactUpload(cidsdk.ArtifactUploadRequest{
					File:   report.Path,
					Module: d.Module.Slug,
					Type:   "report",
					Format: "jacoco",
				})
				if err != nil {
					return err
				}
			}
		}

	} else if d.Module.BuildSystem == string(cidsdk.BuildSystemMaven) {
		mavenWrapper := cidsdk.JoinPath(d.Module.ModuleDir, "mvnw")
		if !a.Sdk.FileExists(mavenWrapper) {
			return fmt.Errorf("maven wrapper not found at %s", mavenWrapper)
		}

		buildArgs := []string{
			`test`,
			`--batch-mode`,
		}
		buildResult, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: fmt.Sprintf("java-exec %s %s", mavenWrapper, strings.Join(buildArgs, " ")),
			WorkDir: d.Module.ModuleDir,
		})
		if err != nil {
			return err
		} else if buildResult.Code != 0 {
			return fmt.Errorf("maven test failed, exit code %d", buildResult.Code)
		}
	}

	return nil
}
