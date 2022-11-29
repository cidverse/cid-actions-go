package techdocs

import (
	"path"

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

	outputDir := path.Join(ctx.Module.ProjectDir, ctx.Module.Slug, "docs")

	if ctx.Module.BuildSystem == string(cidsdk.BuildSystemMkdocs) {
		_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: `techdocs-cli generate --source-dir . --output-dir ` + outputDir + ` --no-docker --etag ${NCI_COMMIT_SHA}`,
			WorkDir: ctx.Module.ModuleDir,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
