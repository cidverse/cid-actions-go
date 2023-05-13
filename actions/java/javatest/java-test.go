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
	MavenVersion string `json:"maven_version"        env:"MAVEN_VERSION"`
}

func (a Action) Execute() (err error) {
	cfg := Config{}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
	}

	// version
	if cfg.MavenVersion == "" {
		cfg.MavenVersion = javacommon.GetVersion(ctx.Env["NCI_COMMIT_REF_TYPE"], ctx.Env["NCI_COMMIT_REF_RELEASE"], ctx.Env["NCI_COMMIT_SHA_SHORT"])
	}

	// test
	if ctx.Module.BuildSystem == string(cidsdk.BuildSystemGradle) {
		gradleWrapper := cidsdk.JoinPath(ctx.Module.ModuleDir, "gradlew")
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
			WorkDir: ctx.Module.ModuleDir,
		})
		if err != nil {
			return err
		} else if testResult.Code != 0 {
			return fmt.Errorf("gradle test failed, exit code %d", testResult.Code)
		}

		// collect and store jacoco test reports
		testReports, err := a.Sdk.FileList(cidsdk.FileRequest{
			Directory:  ctx.Module.ModuleDir,
			Extensions: []string{"jacocoTestReport.xml"},
		})
		for _, report := range testReports {
			if strings.HasSuffix(report.Path, cidsdk.JoinPath("build", "reports", "jacoco", "test", "jacocoTestReport.xml")) {
				err := a.Sdk.ArtifactUpload(cidsdk.ArtifactUploadRequest{
					File:   report.Path,
					Module: ctx.Module.Slug,
					Type:   "report",
					Format: "jacoco",
				})
				if err != nil {
					return err
				}
			}
		}

	} else if ctx.Module.BuildSystem == string(cidsdk.BuildSystemMaven) {
		mavenWrapper := cidsdk.JoinPath(ctx.Module.ModuleDir, "mvnw")
		if !a.Sdk.FileExists(mavenWrapper) {
			return fmt.Errorf("maven wrapper not found at %s", mavenWrapper)
		}

		buildArgs := []string{
			`test`,
			`--batch-mode`,
		}
		buildResult, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: fmt.Sprintf("java-exec %s %s", mavenWrapper, strings.Join(buildArgs, " ")),
			WorkDir: ctx.Module.ModuleDir,
		})
		if err != nil {
			return err
		} else if buildResult.Code != 0 {
			return fmt.Errorf("maven test failed, exit code %d", buildResult.Code)
		}
	}

	return nil
}
