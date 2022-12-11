package java

import (
	"fmt"
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
		var testArgs []string
		testArgs = append(testArgs, fmt.Sprintf(`-Pversion="%s"`, GetVersion(ctx.Env["NCI_COMMIT_REF_TYPE"], ctx.Env["NCI_COMMIT_REF_RELEASE"])))
		testArgs = append(testArgs, `check`)
		testArgs = append(testArgs, `--no-daemon`)
		testArgs = append(testArgs, `--warning-mode=all`)
		testArgs = append(testArgs, `--console=plain`)
		testArgs = append(testArgs, `--stacktrace`)

		_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: fmt.Sprintf("%s %s", javaGradleCmd, strings.Join(testArgs, " ")),
			WorkDir: ctx.Module.ModuleDir,
		})
		if err != nil {
			return err
		}

		// collect test reports

	} else if ctx.Module.BuildSystem == string(cidsdk.BuildSystemMaven) {

	}

	return nil
}
