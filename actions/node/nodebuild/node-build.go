package nodebuild

import (
	"github.com/cidverse/cid-actions-go/actions/node/nodecommon"
	cidsdk "github.com/cidverse/cid-sdk-go"
)

type Action struct {
	Sdk cidsdk.SDKClient
}

type Config struct {
}

func (a Action) Metadata() cidsdk.ActionMetadata {
	return cidsdk.ActionMetadata{
		Name:        "node-build",
		Description: "Builds a node.js project",
		Category:    "build",
		Scope:       cidsdk.ActionScopeModule,
		Rules: []cidsdk.ActionRule{
			{
				Type:       "cel",
				Expression: `MODULE_BUILD_SYSTEM == "npm"`,
			},
		},
		Access: cidsdk.ActionAccess{
			Environment: []cidsdk.ActionAccessEnv{},
			Executables: []cidsdk.ActionAccessExecutable{
				{
					Name: "yarn",
				},
			},
		},
	}
}

func (a Action) Execute() (err error) {
	cfg := Config{}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
	}

	// package.json
	content, err := a.Sdk.FileRead(cidsdk.JoinPath(ctx.Module.ModuleDir, "package.json"))
	if err != nil {
		return err
	}
	pkg, err := nodecommon.ParsePackageJSON(content)
	if err != nil {
		return err
	}

	// build, if script is present
	_, buildScriptPresent := pkg.Scripts[`build`]
	if buildScriptPresent {
		// install
		_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: `yarn install`,
			WorkDir: ctx.Module.ModuleDir,
		})
		if err != nil {
			return err
		}

		// build
		_, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: `yarn build`,
			WorkDir: ctx.Module.ModuleDir,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
