package java

import (
	"fmt"
	"strings"

	cidsdk "github.com/cidverse/cid-sdk-go"
)

type BuildAction struct {
	Sdk cidsdk.SDKClient
}

type BuildConfig struct {
}

func (a BuildAction) Execute() (err error) {
	cfg := BuildConfig{}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
	}

	if ctx.Module.BuildSystem == string(cidsdk.BuildSystemGradle) {
		var buildArgs []string
		buildArgs = append(buildArgs, fmt.Sprintf(`-Pversion="%s"`, GetVersion(ctx.Env["NCI_COMMIT_REF_TYPE"], ctx.Env["NCI_COMMIT_REF_RELEASE"])))
		buildArgs = append(buildArgs, `assemble`)
		buildArgs = append(buildArgs, `--no-daemon`)
		buildArgs = append(buildArgs, `--warning-mode=all`)
		buildArgs = append(buildArgs, `--console=plain`)
		buildArgs = append(buildArgs, `--stacktrace`)

		_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: fmt.Sprintf("%s %s", javaGradleCmd, strings.Join(buildArgs, " ")),
			WorkDir: ctx.Module.ModuleDir,
		})
		if err != nil {
			return err
		}
	} else if ctx.Module.BuildSystem == string(cidsdk.BuildSystemMaven) {

	}

	return nil
}
