package golangbuild

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/cidverse/cid-actions-go/actions/golang/golangcommon"
	"github.com/cidverse/cid-actions-go/util/golang"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/shomali11/parallelizer"
)

type BuildAction struct {
	Sdk cidsdk.SDKClient
}

type BuildConfig struct {
	Platform []golangcommon.Platform `json:"platform"`
}

func (a BuildAction) Execute() error {
	cfg := BuildConfig{}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
	}

	// default to current platform
	if len(cfg.Platform) == 0 {
		// parse platform comment from go.mod
		if len(ctx.Module.Discovery) > 0 && ctx.Module.Discovery[0].File != "" {
			platforms, err := golangcommon.DiscoverPlatformsFromGoMod(ctx.Module.Discovery[0].File)
			if err != nil {
				return err
			}
			cfg.Platform = platforms
		}

		// default to current platform
		if len(cfg.Platform) == 0 {
			cfg.Platform = append(cfg.Platform, golangcommon.Platform{Goos: runtime.GOOS, Goarch: runtime.GOARCH})
		}
	}

	// don't build libraries
	if golang.IsGoLibrary(&ctx.Module) {
		_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "info", Message: "no go files in module root, not attempting to build library projects"})
		return nil
	}

	// build
	if ctx.Module.BuildSystem == "gomod" {
		group := parallelizer.NewGroup()
		defer group.Close()

		for _, p := range cfg.Platform {
			goos := p.Goos
			goarch := p.Goarch
			_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "info", Message: "compile binary", Context: map[string]interface{}{"goos": goos, "goarch": goarch}})

			buildEnv := map[string]string{
				"CGO_ENABLED": "false",
				"GOOS":        goos,
				"GOARCH":      goarch,
				"GOTOOLCHAIN": "local",
			}
			err = group.Add(func() error {
				outputFile := cidsdk.JoinPath(ctx.Config.TempDir, fmt.Sprintf("%s_%s", goos, goarch))

				// build
				buildResult, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
					Command: fmt.Sprintf(`go build -buildvcs=false -ldflags "-s -w -X main.version={NCI_COMMIT_REF_RELEASE} -X main.commit={NCI_COMMIT_HASH} -X main.date={TIMESTAMP_RFC3339} -X main.status={NCI_REPOSITORY_STATUS}" -o %s .`, outputFile),
					WorkDir: ctx.Module.ModuleDir,
					Env:     buildEnv,
				})
				if err != nil {
					return err
				} else if buildResult.Code != 0 {
					return fmt.Errorf("go build failed, exit code %d", buildResult.Code)
				}

				// store result
				err = a.Sdk.ArtifactUpload(cidsdk.ArtifactUploadRequest{
					File:   outputFile,
					Module: ctx.Module.Slug,
					Type:   "binary",
					Format: "go",
				})
				if err != nil {
					return err
				}

				return nil
			})
			if err != nil {
				return err
			}
		}

		err := group.Wait()
		if err != nil {
			return err
		}
	} else {
		return errors.New("build system " + ctx.Module.BuildSystem + " is not supported")
	}

	return nil
}
