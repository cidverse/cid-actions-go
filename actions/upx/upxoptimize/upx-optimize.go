package upxoptimize

import (
	cidsdk "github.com/cidverse/cid-sdk-go"
)

type Action struct {
	Sdk cidsdk.SDKClient
}

type OptimizeConfig struct {
	Mode string `json:"mode"`
}

func (a Action) Metadata() cidsdk.ActionMetadata {
	return cidsdk.ActionMetadata{
		Name: "upx-optimize",
		Description: `
			Optimizes the binary size using upx.

			Notes:
        	- This action requires the UPX_ENABLED environment variable to be set to true.
        	- upx requires a lot of compute resources and runs for a long time, only use this action if you have enough resources available (public ci providers have limits).`,
		Category: "optimize",
		Scope:    cidsdk.ActionScopeProject,
		Rules: []cidsdk.ActionRule{
			{
				Type:       "cel",
				Expression: `ENV["UPX_ENABLED"] == "true"`,
			},
		},
		Access: cidsdk.ActionAccess{
			Environment: []cidsdk.ActionAccessEnv{},
			Executables: []cidsdk.ActionAccessExecutable{
				{
					Name: "upx",
				},
			},
		},
		Input: cidsdk.ActionInput{
			Artifacts: []cidsdk.ActionArtifactType{
				{
					Type: "binary",
				},
			},
		},
	}
}

func (a Action) Execute() (err error) {
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
	files, err := a.Sdk.FileList(cidsdk.FileRequest{Directory: cidsdk.JoinPath(ctx.ProjectDir, "bin")})
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
