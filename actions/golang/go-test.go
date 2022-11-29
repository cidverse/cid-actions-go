package golang

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	cidsdk "github.com/cidverse/cid-sdk-go"
)

type TestAction struct {
	Sdk cidsdk.SDKClient
}

func (a TestAction) Execute() error {
	ctx, err := a.Sdk.ModuleAction(nil)
	if err != nil {
		return err
	}

	covarageDir := path.Join(ctx.Config.ArtifactDir, ctx.Module.Slug, "go-test")
	err = os.MkdirAll(covarageDir, os.ModePerm)
	if err != nil {
		return err
	}
	if ctx.Module.BuildSystem == "gomod" {
		var testArgs []string
		testArgs = append(testArgs, `-vet off`)
		testArgs = append(testArgs, `-cover`)
		testArgs = append(testArgs, `-covermode=count`)
		testArgs = append(testArgs, fmt.Sprintf(`-coverprofile %s/cover.out`, covarageDir))

		if ctx.Config.Debug || ctx.Config.Log["bin-go"] == "debug" {
			testArgs = append(testArgs, `-v`)
		}

		_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "info", Message: "running tests"})
		_, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: fmt.Sprintf("go test %s ./...", strings.Join(testArgs, " ")),
			WorkDir: ctx.Module.ModuleDir,
		})
		if err != nil {
			return errors.New("tests failed: " + err.Error())
		}

		_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "info", Message: "generating coverage report", Context: map[string]interface{}{"output": covarageDir + "/cover.html"}})
		_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: fmt.Sprintf("go tool cover -html %s/cover.out -o %s/cover.html", covarageDir, covarageDir),
			WorkDir: ctx.ProjectDir,
		})
		if err != nil {
			return errors.New("failed to generate html test coverage report: " + err.Error())
		}

	} else {
		return errors.New("build system " + ctx.Module.BuildSystem + " is not supported")
	}

	return nil
}
