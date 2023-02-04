package golang

import (
	"errors"
	"fmt"
	"runtime"

	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/shomali11/parallelizer"
)

type BuildAction struct {
	Sdk cidsdk.SDKClient
}

type BuildConfig struct {
	Platform []Platform `json:"platform"`
}

func (a BuildAction) Execute() error {
	cfg := BuildConfig{}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
	}

	// default to current platform
	if len(cfg.Platform) == 0 {
		cfg.Platform = append(cfg.Platform, Platform{Goos: runtime.GOOS, Goarch: runtime.GOARCH})
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
				"GOPROXY":     "https://goproxy.io,direct",
				"GOOS":        goos,
				"GOARCH":      goarch,
			}

			err := group.Add(func() error {
				buildResult, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
					Command: fmt.Sprintf("go build -buildvcs=false %s-o %s/%s/bin/%s_%s .", ctx.Config.DebugFlag("bin-go", "-v "), ctx.Config.ArtifactDir, ctx.Module.Slug, goos, goarch),
					WorkDir: ctx.Module.ModuleDir,
					Env:     buildEnv,
				})
				if err != nil {
					return err
				} else if buildResult.Code != 0 {
					return fmt.Errorf("go build failed, exit code %d", buildResult.Code)
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
