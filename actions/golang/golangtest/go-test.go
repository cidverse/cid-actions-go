package golangtest

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	cidsdk "github.com/cidverse/cid-sdk-go"
)

type Action struct {
	Sdk cidsdk.SDKClient
}

func (a Action) Metadata() cidsdk.ActionMetadata {
	return cidsdk.ActionMetadata{
		Name:        "go-test",
		Description: "Runs all tests in your go project.",
		Category:    "test",
		Scope:       cidsdk.ActionScopeModule,
		Rules: []cidsdk.ActionRule{
			{
				Type:       "cel",
				Expression: `MODULE_BUILD_SYSTEM == "gomod"`,
			},
		},
		Access: cidsdk.ActionAccess{
			Environment: []cidsdk.ActionAccessEnv{},
			Executables: []cidsdk.ActionAccessExecutable{
				{
					Name: "go",
				},
			},
		},
	}
}

func (a Action) Execute() error {
	ctx, err := a.Sdk.ModuleAction(nil)
	if err != nil {
		return err
	}

	if ctx.Module.BuildSystem != "gomod" {
		return errors.New("build system " + ctx.Module.BuildSystem + " is not supported")
	}

	// paths
	coverageOut := filepath.Join(ctx.Config.TempDir, "cover.out")
	coverageJSON := filepath.Join(ctx.Config.TempDir, "cover.json")
	coverageHTML := filepath.Join(ctx.Config.TempDir, "cover.html")

	// pull dependencies
	pullResult, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: `go get -v -t ./...`,
		WorkDir: ctx.Module.ModuleDir,
		Env: map[string]string{
			"GOTOOLCHAIN": "local",
		},
	})
	if err != nil {
		return err
	} else if pullResult.Code != 0 {
		return fmt.Errorf("go get failed, exit code %d", pullResult.Code)
	}

	// run tests
	testArgs := []string{
		"-vet off",
		"-cover",
		"-covermode=count",
		fmt.Sprintf(`-coverprofile %q`, coverageOut),
		"-parallel=4",
		"-timeout 10s",
	}
	if ctx.Config.Debug || ctx.Config.Log["bin-go"] == "debug" {
		testArgs = append(testArgs, `-v`)
	}
	_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "info", Message: "running tests"})
	testResult, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: fmt.Sprintf("go test %s ./...", strings.Join(testArgs, " ")),
		Env: map[string]string{
			"GOTOOLCHAIN": "local",
		},
		WorkDir: ctx.Module.ModuleDir,
	})
	if err != nil {
		return errors.New("tests failed: " + err.Error())
	} else if testResult.Code != 0 {
		return fmt.Errorf("go test report generation failed, exit code %d", testResult.Code)
	}

	err = a.Sdk.ArtifactUpload(cidsdk.ArtifactUploadRequest{
		Module:        ctx.Module.Slug,
		File:          coverageOut,
		Type:          "report",
		Format:        "go-coverage",
		FormatVersion: "out",
	})
	if err != nil {
		return err
	}

	// json report
	_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "info", Message: "generating json coverage report"})
	coverageJSONResult, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command:       fmt.Sprintf("go test -coverprofile %q -json -covermode=count -parallel=4 -timeout 10s ./...", coverageOut),
		WorkDir:       ctx.Module.ModuleDir,
		CaptureOutput: true,
	})
	if err != nil {
		return errors.New("failed to generate json test coverage report: " + err.Error())
	} else if coverageJSONResult.Code != 0 {
		return fmt.Errorf("go test report generation failed, exit code %d", coverageJSONResult.Code)
	}

	err = a.Sdk.FileWrite(coverageJSON, []byte(coverageJSONResult.Stdout))
	if err != nil {
		return errors.New("failed to store json test coverage report on filesystem: " + err.Error())
	}

	err = a.Sdk.ArtifactUpload(cidsdk.ArtifactUploadRequest{
		Module:        ctx.Module.Slug,
		File:          coverageJSON,
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
		Command: fmt.Sprintf("go tool cover -html %q -o %q", coverageOut, coverageHTML),
		WorkDir: ctx.ProjectDir,
	})
	if err != nil {
		return errors.New("failed to generate html test coverage report: " + err.Error())
	}

	err = a.Sdk.ArtifactUpload(cidsdk.ArtifactUploadRequest{
		Module:        ctx.Module.Slug,
		File:          coverageHTML,
		Type:          "report",
		Format:        "go-coverage",
		FormatVersion: "html",
	})
	if err != nil {
		return err
	}

	return nil
}
