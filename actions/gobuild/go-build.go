package gobuild

import (
	"errors"
	"fmt"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/shomali11/parallelizer"
	"runtime"
)

type Action struct {
	Sdk cidsdk.SDKClient
}

type Config struct {
	Platform []Platform `json:"platform"`
}

type Platform struct {
	Goos   string `required:"true" json:"goos"`
	Goarch string `required:"true" json:"goarch"`
}

func (a Action) Execute() error {
	cfg := Config{}
	ctx, err := a.Sdk.PrepareAction(&cfg)
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
			_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "info", Message: "go-build", Context: map[string]interface{}{"goos": goos, "goarch": goarch}})

			buildEnv := map[string]string{
				"CGO_ENABLED": "false",
				"GOPROXY":     "https://goproxy.io,direct",
				"GOOS":        goos,
				"GOARCH":      goarch,
			}

			err := group.Add(func() error {
				_, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
					Command: fmt.Sprintf("go build -buildvcs=false %s-o %s/%s/bin/%s_%s .", ctx.Config.DebugFlag("bin-go", "-v "), ctx.Config.ArtifactDir, ctx.Module.Slug, goos, goarch),
					WorkDir: ctx.Module.ModuleDir,
					Env:     buildEnv,
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
