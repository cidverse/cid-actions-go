package upx

import (
	"path"

	cidsdk "github.com/cidverse/cid-sdk-go"
)

type OptimizeAction struct {
	Sdk cidsdk.SDKClient
}

type OptimizeConfig struct {
	Mode string `json:"mode"`
}

func (a OptimizeAction) Execute() (err error) {
	cfg := OptimizeConfig{}
	ctx, err := a.Sdk.ProjectAction(&cfg)
	if err != nil {
		return err
	}

	// mode
	if cfg.Mode == "" {
		cfg.Mode = "--lzma"
	}

	// files
	files, err := a.Sdk.FileList(cidsdk.FileRequest{Directory: path.Join(ctx.ProjectDir, "bin")})
	if err != nil {
		return err
	}
	for _, file := range files {
		_, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: `upx ` + cfg.Mode + ` ` + file.Path,
			WorkDir: ctx.ProjectDir,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
