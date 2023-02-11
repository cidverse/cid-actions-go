package java

import (
	"fmt"
	"path"
	"strings"

	cidsdk "github.com/cidverse/cid-sdk-go"
)

type TestAction struct {
	Sdk cidsdk.SDKClient
}

type TestConfig struct {
}

func (a TestAction) Execute() (err error) {
	cfg := TestConfig{}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
	}

	if ctx.Module.BuildSystem == string(cidsdk.BuildSystemGradle) {
		testArgs := []string{
			fmt.Sprintf(`-Pversion=%q`, GetVersion(ctx.Env["NCI_COMMIT_REF_TYPE"], ctx.Env["NCI_COMMIT_REF_RELEASE"])),
			`check`,
			`--no-daemon`,
			`--warning-mode=all`,
			`--console=plain`,
			`--stacktrace`,
		}
		testResult, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: fmt.Sprintf("%s %s", javaGradleCmd, strings.Join(testArgs, " ")),
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
			if strings.HasSuffix(report.Path, path.Join("build", "reports", "jacoco", "test", "jacocoTestReport.xml")) {
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

	}

	return nil
}
