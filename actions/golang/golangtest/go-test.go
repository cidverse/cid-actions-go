package golangtest

import (
	"errors"
	"fmt"
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

	if ctx.Module.BuildSystem != "gomod" {
		return errors.New("build system " + ctx.Module.BuildSystem + " is not supported")
	}

	// paths
	coverageDir := ctx.Config.TempDir

	// run tests
	testArgs := []string{
		"-vet off",
		"-cover",
		"-covermode=count",
		fmt.Sprintf(`-coverprofile %s/cover.out`, coverageDir),
		"-parallel=4",
		"-timeout 10s",
	}
	if ctx.Config.Debug || ctx.Config.Log["bin-go"] == "debug" {
		testArgs = append(testArgs, `-v`)
	}
	_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "info", Message: "running tests"})
	testResult, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: fmt.Sprintf("go test %s ./...", strings.Join(testArgs, " ")),
		WorkDir: ctx.Module.ModuleDir,
	})
	if err != nil {
		return errors.New("tests failed: " + err.Error())
	} else if testResult.Code != 0 {
		return fmt.Errorf("go test report generation failed, exit code %d", testResult.Code)
	}

	err = a.Sdk.ArtifactUpload(cidsdk.ArtifactUploadRequest{
		Module:        ctx.Module.Slug,
		File:          coverageDir + "/cover.out",
		Type:          "report",
		Format:        "go-coverage",
		FormatVersion: "out",
	})
	if err != nil {
		return err
	}

	// json report
	_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "info", Message: "generating json coverage report"})
	coverageJsonResult, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command:       fmt.Sprintf("go test -coverprofile %s/cover.out -json -covermode=count -parallel=4 -timeout 10s ./...", coverageDir),
		WorkDir:       ctx.Module.ModuleDir,
		CaptureOutput: true,
	})
	if err != nil {
		return errors.New("failed to generate json test coverage report: " + err.Error())
	} else if coverageJsonResult.Code != 0 {
		return fmt.Errorf("go test report generation failed, exit code %d", coverageJsonResult.Code)
	}

	err = a.Sdk.FileWrite(coverageDir+"/cover.json", []byte(coverageJsonResult.Stdout))
	if err != nil {
		return errors.New("failed to store json test coverage report on filesystem: " + err.Error())
	}

	err = a.Sdk.ArtifactUpload(cidsdk.ArtifactUploadRequest{
		Module:        ctx.Module.Slug,
		File:          coverageDir + "/cover.json",
		Type:          "report",
		Format:        "go-coverage",
		FormatVersion: "json",
	})
	if err != nil {
		return err
	}

	// html report
	_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "info", Message: "generating html coverage report"})
	_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: fmt.Sprintf("go tool cover -html %s/cover.out -o %s/cover.html", coverageDir, coverageDir),
		WorkDir: ctx.ProjectDir,
	})
	if err != nil {
		return errors.New("failed to generate html test coverage report: " + err.Error())
	}

	err = a.Sdk.ArtifactUpload(cidsdk.ArtifactUploadRequest{
		Module:        ctx.Module.Slug,
		File:          coverageDir + "/cover.html",
		Type:          "report",
		Format:        "go-coverage",
		FormatVersion: "html",
	})
	if err != nil {
		return err
	}

	return nil
}
